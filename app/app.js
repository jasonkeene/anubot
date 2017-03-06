'use strict';

const electron = require('electron'),
      child_process = require('child_process'),
      electron_debug = require('electron-debug'),
      electron_settings = require('electron-settings');

// setup settings
electron_settings.configure({
    prettify: true,
});
electron_settings.defaults({
    api: 'wss://api.anubot.io/api',
    origin: 'https://api.anubot.io',
});
electron_settings.applyDefaultsSync();

// enable dev tools
electron_debug({
    showDevTools: true,
});

var mainWindow = null; // keep main object from being GC'd

electron.app.on('ready', () => {
    var windowOpts = {
        width: 560,
        height: 620,
        minWidth: 500,
        minHeight: 400,
        frame: false,
        backgroundColor: '#24272f',
        show: false,
        resizable: false,
    };
    mainWindow = new electron.BrowserWindow(windowOpts);
    mainWindow.loadURL('file://' + __dirname + '/app.html');
    mainWindow.once('ready-to-show', () => {
      mainWindow.show();
    })
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
