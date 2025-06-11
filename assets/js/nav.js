function toggleMobileNav() {
	const sidebar = document.getElementById('sidebar');
	const overlay = document.getElementById('sidebar-overlay');

	sidebar.classList.toggle('-translate-x-full');
	overlay.classList.toggle('hidden');
}

function closeMobileNav() {
	const sidebar = document.getElementById('sidebar');
	const overlay = document.getElementById('sidebar-overlay');

	sidebar.classList.add('-translate-x-full');
	overlay.classList.add('hidden');
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
