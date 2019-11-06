# publish-q-sys-plugin

GitHub action to publish Q-Sys Plugins to Nuget server (Built for <https://github.com/q-sys-community>)

## Inputs

### `nuget-host`

**Required** The host for the NuGet server in format `"http://host/path/"`

### `nuspec-file`

The name of the .nuxped file, defaults to first found in root.

### `api-key`

The API Key that grants write access to the NuGet server

### `version`

The version to tag this release as. Defaults to 'commit-tag' which will try to extract it from the commit tag.

## Outputs

None (for now)

## Example usage

uses: soloworks/publish-q-sys-plugin
with:
  nuget-host: 'http://myhost/mypluginpath/'
  api-key: ${{ secrets.NUGET_APIKEY }}