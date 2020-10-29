using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Azure.KeyVault;
using Microsoft.Azure.Services.AppAuthentication;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Configuration.AzureKeyVault;
using Microsoft.Extensions.Logging;

namespace incrementally_backend
{
    public class Program
    {
        public static void Main(string[] args)
        {
            BuildWebHost(args).Run();
        }

        public static IWebHost BuildWebHost(string[] args) =>
        WebHost.CreateDefaultBuilder(args)
        .ConfigureAppConfiguration((ctx, builder) =>
        {
            var keyVaultEndpoint = GetKeyVaultEndpoint();
            if (!string.IsNullOrEmpty(keyVaultEndpoint))
            {
                builder
                    .AddJsonFile("appsettings.json", optional: true, reloadOnChange: true);
                if (Environment.GetEnvironmentVariable("ASPNETCORE_ENVIRONMENT") == "Development")
                {
                    // If this is throwing an exception it's probably because you need to create secrets.json.
                    // See README for more info.
                    builder.AddJsonFile($"appsettings.Development.json", optional: true, reloadOnChange: true);
                    builder.AddUserSecrets<Program>();
                } else
                {
                    var azureServiceTokenProvider = new AzureServiceTokenProvider();
                    var keyVaultClient = new KeyVaultClient(
                        new KeyVaultClient.AuthenticationCallback(azureServiceTokenProvider.KeyVaultTokenCallback));
                    builder
                        .AddAzureKeyVault(keyVaultEndpoint, keyVaultClient, new DefaultKeyVaultSecretManager());
                }
            }
        })
        .UseStartup<Startup>()
        .Build();

        private static string GetKeyVaultEndpoint() => "https://incrementally.vault.azure.net";
    }
}
