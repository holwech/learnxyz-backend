using Newtonsoft.Json;

namespace Database.Models
{
  public class RecordingEntry {
    [JsonProperty(PropertyName = "id")]
    public string Id { get; set; }
    [JsonProperty(PropertyName = "recording")]
    public string Recording { get; set; }
  }
}