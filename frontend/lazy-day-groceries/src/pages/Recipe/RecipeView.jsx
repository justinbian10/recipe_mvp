import { Fragment, useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import ApiClient from '/src/ApiClient.js';
import './RecipeView.css';

function RecipeView() {
	const [recipe, setRecipe] = useState({});
	const { recipeId } = useParams();

	useEffect(() => {
		async function startFetching() {
			const apiClient = new ApiClient("http://localhost:8080/v1");
			const recipeObj = await apiClient.getRecipe(recipeId);
			console.log(recipeObj);
			setRecipe(recipeObj.recipe);
		}
		startFetching();
	}, [])

	return (
		<div className='recipe-view-main'>
			<RecipeHeader recipe={recipe} />
			<RecipeBody recipe={recipe} />
		</div>
	)
}

function RecipeHeader({ recipe }) {
	return (
		<div className='recipe-view-header-container'>
			<img className='recipe-view-picture' src={recipe.image_url} alt='recipe photograph' />
			<h1 className='recipe-view-title'>{recipe.title}</h1> 
			<button className='recipe-view-button edit'>Edit</button>
			<button className='recipe-view-button delete'>Delete</button>
			<p className='recipe-view-desc'>{recipe.description}</p>
			<div className='recipe-view-info-container'>
				<h3 className='info servings'>{recipe.servings} Servings</h3>
				<h3 className='info cooktime'>{recipe.cooktime} Minutes</h3>
			</div>
		</div>
	)
}

function createIngredientsDisplay(ingredients) {
	if (ingredients) {
		return ingredients.map((ingredient, index) =>
			<div className='recipe-view-ingredient' key={index}>
				<li className={'recipe-view-ingredient-name ingredient-name-' + index} >
					{ingredient.name}
				</li>
				<li className={'recipe-view-ingredient-amount ingredient-amount-' + index} >
					{`${ingredient.amount} ${ingredient.unit}`}
				</li>
			</div>)
	}
	return []
}

function createStepsDisplay(steps) {
	if (steps) {
		return steps.map((step, index) =>
			<li className={'recipe-view-step step-' + index} key={index} >
				<p className='recipe-view-step-desc'>{step.description}</p>
			</li>)
	}
	return []
}


function RecipeBody({ recipe }) {
	return (
		<div className='recipe-view-body-container'>
			<h2 className='ingredient-label label'>Ingredients</h2>
			<ul className='ingredient-container'>
				{createIngredientsDisplay(recipe.ingredients)}
			</ul>
			<h2 className='step-label label'>Recipe Steps</h2>
			<ol className='step-container'>
				{createStepsDisplay(recipe.steps)}
			</ol>
		</div>
	)
}

export default RecipeView
