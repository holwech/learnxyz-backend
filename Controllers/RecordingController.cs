using System;
using System.Collections.Generic;
using System.Security.Claims;
using System.Threading.Tasks;
using Database.Models;
using incrementally_backend.Application;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;

namespace incrementally.Controllers
{
    [Route("api")]
    [ApiController]
    [Authorize]
    public class RecordingController : ControllerBase
    {
        private readonly RecordingHandler _recordingHandler;

        public RecordingController(RecordingHandler recordingHandler)
        {
            _recordingHandler = recordingHandler;
        }

        [HttpGet]
        [AllowAnonymous]
        public string Get()
        {
            return "Server is running";
        }

        [HttpGet]
        [AllowAnonymous]
        [Route("recording/{id?}")]
        public async Task<RecordingEntry> Recording(string id)
        {
            return await _recordingHandler.GetRecording(id);
        }

        [HttpGet]
        [AllowAnonymous]
        [Route("metadata/{id}")]
        public async Task<RecordingMetadata> SingleMetadata(string id)
        {
            return await _recordingHandler.GetMetadata(id);
        }

        [HttpGet]
        [AllowAnonymous]
        [Route("metadata")]
        public async Task<IEnumerable<RecordingMetadata>> Metadata()
        {
            return await _recordingHandler.GetTopMetadata();
        }


        [HttpPost]
        [Route("create")]
        public async Task<RecordingMetadata> CreateAsync(UserRecordingInput data)
        {
            var recordingMetadata = new RecordingMetadata {
                Title = data.Title,
                Description = data.Description,
                CreatedBy = User.FindFirstValue(ClaimTypes.NameIdentifier),
                GivenName = User.FindFirstValue(ClaimTypes.GivenName),
                Surname = User.FindFirstValue(ClaimTypes.Surname),
                Id = Guid.NewGuid().ToString(),
                CreatedAt = DateTimeOffset.Now
            };
            await _recordingHandler.Add(recordingMetadata, data.Recording);
            return recordingMetadata;
        }
    }

    public class UserRecordingInput
    {
        public string Recording { get; set; }
        public string Title { get; set; }
        public string Description { get; set; }
    }
}
