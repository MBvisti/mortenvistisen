function changeTheme(status) {
    const theme = document.querySelector("[data-theme]")

    const lightElement = document.querySelector("[data-light-theme]")
    const darkElement = document.querySelector("[data-dark-theme]")

    switch (status.checked) {
        case true:
            lightElement.style.visibility = "hidden";
            darkElement.style.visibility = "visible";

            theme.setAttribute("data-theme", "dracula");

            localStorage.setItem('theme', 'dark');

            break;
        case false:
            lightElement.style.visibility = "visible";
            darkElement.style.visibility = "hidden";

            theme.setAttribute("data-theme", "garden");

            localStorage.setItem('theme', 'light');

            break;
        default:
            lightElement.style.visibility = "hidden";
            darkElement.style.visibility = "visible";

            theme.setAttribute("data-theme", "darcula");

            localStorage.setItem('theme', 'dark');
    }
}

function setTheme() {
    const theme = document.querySelector("[data-theme]")
    const currentTheme = localStorage.getItem('theme');
    const lightElement = document.querySelector("[data-light-theme]")
    const darkElement = document.querySelector("[data-dark-theme]")
    const themeSelector = document.querySelector("[data-theme-selector]")

    switch (currentTheme) {
        case "dark":
            theme.setAttribute("data-theme", "dracula");

            break;
        case "light":

            theme.setAttribute("data-theme", "garden");

            break;
        default:
            theme.setAttribute("data-theme", "darcula");
    }
}

function setThemeAfterLoad() {
    const theme = document.querySelector("[data-theme]")
    const currentTheme = localStorage.getItem('theme');
    const lightElement = document.querySelector("[data-light-theme]")
    const darkElement = document.querySelector("[data-dark-theme]")
    const themeSelector = document.querySelector("[data-theme-selector]")

    switch (currentTheme) {
        case "dark":
            theme.setAttribute("data-theme", "dracula");

            break;
        case "light":
            lightElement.style.visibility = "visible";
            darkElement.style.visibility = "hidden";

            themeSelector.removeAttribute("checked");

            theme.setAttribute("data-theme", "garden");

            break;
        default:
            theme.setAttribute("data-theme", "darcula");
    }
}
