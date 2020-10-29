using Azure;
using Azure.Storage.Blobs;
using Azure.Storage.Blobs.Models;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;

namespace incrementally_backend.Services
{
    public class AzureBlobStorageConnector : IStorageConnector
    {
        private Dictionary<string, BlobContainerClient> _containers = new Dictionary<string, BlobContainerClient>();
        public void Initialize(string connectionString, List<string> containerNames) {
            var client = new BlobServiceClient(connectionString);
            var tasks = new List<(string, Task<BlobContainerClient>)>();
            containerNames.ForEach(containerName => tasks.Add((containerName, CreateContainer(client, containerName))));
            foreach (var task in tasks)
            {
                try
                {
                    _containers.Add(task.Item1, task.Item2.GetAwaiter().GetResult());
                } catch (RequestFailedException)
                {
                    _containers.Add(task.Item1, client.GetBlobContainerClient(task.Item1));
                }
            }
        }

        public async Task UploadAsync(string fileName, string fileContent, string containerName)
        {
            var container = _containers[containerName];
            BlobClient blobClient = container.GetBlobClient(fileName);
            await blobClient.UploadAsync(ToStream(fileContent));
        }

        public async Task<string> DownloadAsync(string fileName, string containerName)
        {
            var container = _containers[containerName];
            BlobClient blobClient = container.GetBlobClient(fileName);
            BlobDownloadInfo download = await blobClient.DownloadAsync();
            MemoryStream stream = new MemoryStream();
            await download.Content.CopyToAsync(stream);
            // Could this potentially cause memory issues? Should be looked into at some point.
            return System.Text.Encoding.UTF8.GetString(stream.ToArray());
        }

        private async Task<BlobContainerClient> CreateContainer(BlobServiceClient blobClient, string containerName)
        {
            return await blobClient.CreateBlobContainerAsync(containerName);
        }

        private Stream ToStream(string file)
        {
            var stream = new MemoryStream();
            var writer = new StreamWriter(stream);
            writer.Write(file);
            writer.Flush();
            stream.Position = 0;
            return stream;
        }
    }
}
