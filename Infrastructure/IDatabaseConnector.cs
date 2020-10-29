using Database.Models;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace incrementally.Services
{
    public interface IDatabaseConnector
    {
        Task Initialize(List<string> containerNames, string account, string key);
        Task CreateContainer(string containerName);
        Task DeleteRecording(string id);
        Task AddRecording(RecordingMetadata recordingMetadata);
        Task<IEnumerable<RecordingMetadata>> GetTopRecordings();
        Task<RecordingMetadata> GetRecording(string id);
    }
}