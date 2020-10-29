using System;
using Newtonsoft.Json;

namespace Database.Models
{
  public class RecordingMetadata {
    [JsonProperty(PropertyName = "createdBy")]
    public string CreatedBy { get; set; }

    [JsonProperty(PropertyName = "createdAt")]
    public DateTimeOffset CreatedAt { get; set; }

    [JsonProperty(PropertyName = "givenName")]
    public string GivenName { get; set; }

    [JsonProperty(PropertyName = "surname")]
    public string Surname { get; set; }
    
    [JsonProperty(PropertyName = "id")]
    public string Id { get; set; }

    [JsonProperty(PropertyName = "title")]
    public string Title { get; set; }

    [JsonProperty(PropertyName = "description")]
    public string Description { get; set; }
  }
}