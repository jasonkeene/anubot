'use strict';

const electron = require('electron');
const child_process = require('child_process')

var mainWindow = null; // keep main object from being GC'd

electron.app.on('ready', () => {
    var go_proc = child_process.spawn("./anubot-server")

    go_proc.stdout.on('data', (data) => {
        console.log(`stdout: ${data}`);
    });
    go_proc.stderr.on('data', (data) => {
        console.log(`stderr: ${data}`);
    });
    go_proc.on('close', (code) => {
        console.log(`child process exited with code ${code}`);
    });

    mainWindow = new electron.BrowserWindow({
        width: 1024,
        height: 768,
    });
    mainWindow.loadURL('file://' + __dirname + '/index.html');
    mainWindow.on('closed', () => {
        mainWindow = null;
        go_proc.kill()
    });
    mainWindow.webContents.openDevTools();
});

// tear down app when windows are closed
electron.app.on('window-all-closed', () => {
    // on osx processes typically run until the user hits cmd+q
    if (process.platform != 'darwin') {
        electron.app.quit();
    }
});
