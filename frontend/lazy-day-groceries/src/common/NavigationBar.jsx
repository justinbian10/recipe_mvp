import "./NavigationBar.css";
function NavigationBar() {
	return (
		<nav>
			<ul>
				<div className="left">
					<li>
						<a href="/">
							<img src="/src/assets/images/logo.jpg" className="logo" alt="website logo" />
							<h2 className="logo-text">LDG</h2>
						</a>
					</li>
					<li>
						<a href="/loadout">Loadout</a>
					</li>
					<li>
						<a href="/manage-recipes">Manage Recipes</a>
					</li>
				</div>
				<div className="right">
					<li>
						<a href="/login">Login</a>
					</li>
				</div>
			</ul>
		</nav>
	)
}

export default NavigationBar
