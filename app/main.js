'use strict';

const electron = require('electron'),
      child_process = require('child_process'),
      electron_debug = require('electron-debug');

electron_debug({showDevTools: true});

var mainWindow = null; // keep main object from being GC'd

electron.app.on('ready', () => {
    var windowOpts = {
        width: 1024,
        height: 768,
        frame: false,
    };
    mainWindow = new electron.BrowserWindow(windowOpts);
    mainWindow.loadURL('file://' + __dirname + '/index.html');
    mainWindow.on('closed', () => {
        mainWindow = null;
    });
});

// tear down app when windows are closed
electron.app.on('window-all-closed', () => {
    // on osx processes typically run until the user hits cmd+q
    if (process.platform !== 'darwin') {
        electron.app.quit();
    }
});
