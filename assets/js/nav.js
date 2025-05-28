function toggleMobileNav() {
	const sidebar = document.getElementById('sidebar');
	const overlay = document.querySelector('.sidebar-overlay');

	sidebar.classList.toggle('open');
	overlay.classList.toggle('open');
}

function closeMobileNav() {
	const sidebar = document.getElementById('sidebar');
	const overlay = document.querySelector('.sidebar-overlay');

	sidebar.classList.remove('open');
	overlay.classList.remove('open');
}

document.addEventListener('DOMContentLoaded', function() {
	const navToggle = document.querySelector('.nav-toggle');
	const navMenu = document.querySelector('.nav-menu');

	if (navToggle && navMenu) {
		navToggle.addEventListener('click', () => {
			navMenu.classList.toggle('nav-menu-active');
			navToggle.classList.toggle('nav-toggle-active');
		});
	}
});
