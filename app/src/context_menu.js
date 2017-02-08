
const {remote} = require('electron');
const {Menu, MenuItem} = remote;

function register() {
    const userMenu = new Menu();
    userMenu.append(new MenuItem({
        label: 'MenuItem1',
        click() {
            console.log("item 1 clicked for username:", currentUsername);
        },
    }))
    userMenu.append(new MenuItem({
        label: 'MenuItem2',
        click() {
            console.log("item 2 clicked for username:", currentUsername);
        },
    }));

    var currentUsername = null;
    window.addEventListener('contextmenu', (e) => {
        e.preventDefault();
        if (e.target.className === "nick") {
            currentUsername = e.target.textContent;
            userMenu.popup(remote.getCurrentWindow())
        }
    }, false);
}

module.exports = {
    register,
};
