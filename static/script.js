// Pico respects the user's system-wide dark mode preference, but we can also allow the user to manually toggle dark mode on or off.
const getPreferredScheme = () => window?.matchMedia?.('(prefers-color-scheme:dark)')?.matches ? 'dark' : 'light';

const setTheme = (theme) => {
    document.documentElement.setAttribute('data-theme', theme);
}

const toggleTheme = () => {
    const currentTheme = document.documentElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    localStorage.setItem('theme', newTheme);
    setTheme(newTheme);
}

document.addEventListener('DOMContentLoaded', () => {
    const darkModeButton = document.getElementById('dark-mode');
    const preferredScheme = localStorage.getItem('theme') || getPreferredScheme();

    setTheme(preferredScheme);
    if (preferredScheme === 'dark') {
        darkModeButton.checked = true;
    }
    if (darkModeButton) {
        darkModeButton.addEventListener('click', toggleTheme);
    }
});
