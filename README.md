### Configure project
1. Open `sfdx-project.json`
2. Make current path in `packageDirectories` as not default (set `false`)
3. Add property `watch` and set `true`
4. Add new path to `packageDirectories` in `sfdx-project.json` and make it default (set `true`)
5. Add new path to `.gitignore` and `.forceignore` (but it is not required)

Example

```
{
  "packageDirectories": [
    {
      "path": "force-app",
      "default": false,
      "watch": true
    },{
      "path": "force-app-build",
      "default": true
    }
  ],
  "namespace": "",
  "sfdcLoginUrl": "https://login.salesforce.com",
  "sourceApiVersion": "51.0"
}
```

### Configure WebStorm
1. Put `lwc_watcher_(mac|linux|win)` to project root
2. Open `Preferences -> File Watcher`
3. Click button `Add`
4. Put any name
5. File Type -> Any
6. Select bin file in `Program` field
7. Put `-f $FilePath$` to `Arguments` field
8. If watcher is not in root project dir, put argument `-project-dir $ProjectFileDir$`
9. For better performance you can uncheck `Auto-check edited files to trigger the watcher` and run watcher only if press `Ctrl + S` (`Cmd + S` on MaC)
10. Click `Ok` and `Applay`

If you have lwc components already, just run `./lwc_watcher_(mac|linux|win) --first` in CLI for prepare all yours components
