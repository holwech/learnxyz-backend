using System.Collections.Generic;
using System.Threading.Tasks;

namespace incrementally_backend.Services
{
    public interface IStorageConnector
    {
        public void Initialize(string connectionString, List<string> containerNames);
        public Task<string> DownloadAsync(string fileName, string containerName);
        public Task UploadAsync(string fileName, string fileContent, string containerName);
    }
}
