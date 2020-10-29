using Database.Models;
using incrementally.Services;
using incrementally_backend.Services;
using System;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace incrementally_backend.Application
{
    public class RecordingHandler
    {
        private readonly IDatabaseConnector _database;
        private readonly IStorageConnector _storage;

        public RecordingHandler(IDatabaseConnector database, IStorageConnector storage)
        {
            _database = database;
            _storage = storage;
        }

        public async Task Add(RecordingMetadata recordingMetadata, string recordingData)
        {
            var task = _database.AddRecording(recordingMetadata);
            await _storage.UploadAsync(recordingMetadata.Id, recordingData, "recordings");
            await task;
        }

        public async Task<RecordingMetadata> GetMetadata(string id)
        {
            return await _database.GetRecording(id);
        }

        public async Task<IEnumerable<RecordingMetadata>> GetTopMetadata()
        {
            return await _database.GetTopRecordings();
        }

        public async Task<RecordingEntry> GetRecording(string id)
        {
            return new RecordingEntry
            {
                Id = id,
                Recording = await _storage.DownloadAsync(id, "recordings")
            };
        }

        public async Task<RecordingMetadata> Delete(string userId, string entryId)
        {
            var entry = await _database.GetRecording(entryId);
            if (userId == entry.CreatedBy)
            {
                await _database.DeleteRecording(entryId).ConfigureAwait(false);
                return entry;
            } else
            {
                throw new UnauthorizedAccessException("User not authorized to delete this item");
            }
        }
    }
}