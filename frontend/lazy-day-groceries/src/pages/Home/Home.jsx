import './Home.css';
import ApiClient from '/src/ApiClient.js';
import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom'

function Home() {
	return (
		<main>
			<h1>Lazy Day Groceries</h1>
			<p>Find a recipe!</p>
			<SearchBar />
			<div className='recipe-preview'>
				<h3>Like These!</h3>
				<RecipePreview />
			</div>
		</main>
	)
}

function SearchBar() {
	return (
		<form action='/search'>
			<input type='text' name='search' id='search' placeholder='🔎' />
		</form>
	)	
}

function RecipePreview() {
	const [recipes, setRecipes] = useState([]);
	useEffect(() => {
		async function startFetching() {
			const apiClient = new ApiClient('http://localhost:8080/v1');
			const recipesObj = await apiClient.getAllRecipes();
			setRecipes(recipesObj.recipes);
		}
		startFetching();
	}, []);
	async function addExampleRecipe(recipe) {
		return ExampleRecipe(recipe.id, recipe.image_url, recipe.title);	
	}
	return (
		<div className='recipe-preview-layout'>
			{recipes.map((recipe, index) => {
				return <ExampleRecipe key={recipe.id} recipeId={recipe.id} imageUrl={recipe.image_url} title={recipe.title} />
			})}
		</div>
	)
}

function ExampleRecipe({ recipeId, imageUrl, title }) {
	return (
		<Link to={'/recipe/'+recipeId}>
			<div className='recipe-card'>
				<img className='home-preview-picture' src={imageUrl} />
				<h2 className='home-title'>{title}</h2>
			</div>
		</Link>
	)
}

export default Home
