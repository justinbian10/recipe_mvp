import { useState, useEffect, useRef } from 'react';
import { Link } from 'react-router-dom'
import "./Loadout.css";
import ApiClient from "/src/ApiClient.js";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faMinus, faXmark } from '@fortawesome/free-solid-svg-icons';

function Loadout() {
	const [pantryQuery, setPantryQuery] = useState('');
	const pantryQueryRef = useRef(null);
	const [ingredients, setIngredients] = useState([]);
	const [isIngredientsOpen, setIsIngredientsOpen] = useState(false);
	const [currentIngredients, setCurrentIngredients] = useState({});

	const [recipeQuery, setRecipeQuery] = useState('');
	const [recipes, setRecipes] = useState({});
	const [loadout, setLoadout] = useState({});

	useEffect(() => {
		const apiClient = new ApiClient('http://localhost:8080/v1');
		const ingredients = apiClient.getAllIngredients('');
		const recipes = apiClient.getAllRecipes();
		Promise.all([ingredients, recipes]).then((values) => {
			setIngredients(values[0].ingredients);
			const nextRecipes = {};
			values[1].recipes.forEach((recipe, index) => {
				nextRecipes[recipe.id] = recipe;
			});
			setRecipes(nextRecipes);
		});
	}, []);

	useEffect(() => {
		const handleClickOutside = (event) => {
			if (pantryQueryRef.current && !pantryQueryRef.current.contains(event.target)) {
				setIsIngredientsOpen(false);
			} else {
				setIsIngredientsOpen(true);
			}
		};

		document.addEventListener('mousedown', handleClickOutside);

		return () => document.removeEventListener('mousedown', handleClickOutside);
	}, [pantryQueryRef]);


	function handleIngredientClick(event, ingredient) {
		if (!(ingredient in currentIngredients)) {
			setCurrentIngredients({
				...currentIngredients,
				[ingredient]: 1
			});
		}
	}

	function handleIncrementIngredientClick(event, ingredient, multiple) {
		const { [ingredient]: _, ...remainingIngredients } = currentIngredients;
		const nextAmount = currentIngredients[ingredient] + (1 * multiple);
		if (nextAmount > 0) {
			setCurrentIngredients({
				...currentIngredients,
				[ingredient]: nextAmount
			});
		} else {
			setCurrentIngredients(remainingIngredients);
		}
	}

	function renderQueryIngredients() {
		const queryIngredients = ingredients.filter((ingredient, index) => ingredient.toLowerCase().includes(pantryQuery.toLowerCase()));
		return isIngredientsOpen && queryIngredients && (<ul className="query-ingredients">
			{queryIngredients.map((ingredient, index) => (
				<li className='query-ingredient' key={index} onClick={e => handleIngredientClick(e, ingredient)}>{ingredient}</li>
			))}
		</ul>);
	}

	function renderQueryRecipes() {
		const queryRecipes = Object.entries(recipes).filter(([recipeId, recipe], index) => recipe.title.toLowerCase().includes(recipeQuery.toLowerCase()));
		return queryRecipes && queryRecipes.map(([recipeId, recipe], index) => (
			<tr className={(recipeId in loadout) && 'in-use'}>
				<td><Link className='recipe-link' target="_blank" to={`/recipe/${recipeId}`}>{recipe.title}</Link></td>
				<td>{recipe.servings}</td>
				<td>{recipe.cooktime}</td>
				<td className='add-button-cell'><button className='add-button' onClick={e => handleAddRecipeClick(e, recipe)}><FontAwesomeIcon icon={faPlus} /></button></td>
			</tr>));
	}

	function handleAddRecipeClick(event, recipe) {
		if (!(recipe.id in loadout)) {
			setLoadout({
				...loadout,
				[recipe.id]: 1
			});
		}
	}

	function handleIncrementRecipeClick(event, recipeId, multiple) {
		const { [recipeId]: _, ...remainingRecipeIds } = loadout;
		const nextAmount = loadout[recipeId] + (1 * multiple);
		if (nextAmount > 0) {
			setLoadout({
				...loadout,
				[recipeId]: nextAmount
			});
		} else {
			setLoadout(remainingRecipeIds);
		}
	}
	console.log(currentIngredients);

	return (
		<div className='loadout-container'>
			<div className='current-ingredients-container'>
				<h2 className='loadout-ingredients-label'>Current Pantry {Object.keys(currentIngredients).length > 0 && `(${Object.keys(currentIngredients).length})`}</h2>
				<div className='current-ingredients'>
					<div className='ingredient-query-container' ref={pantryQueryRef}>
						<input type='search' id='ingredient-query' name='ingredient-query' 
							value={pantryQuery}
							onChange={e => setPantryQuery(e.target.value)}
						/>
						{renderQueryIngredients()}
					</div>
					<ul className='added-ingredients'>
						{Object.entries(currentIngredients).map(([ingredient, amount], index) => (
							<li className='added-ingredient' key={index}>
								<p className='loadout-ingredient-name'>{ingredient}</p>
								<button className='decrease-ingredient' onClick={e => handleIncrementIngredientClick(e, ingredient, -1)}>{currentIngredients[ingredient] > 1 ? <FontAwesomeIcon icon={faMinus} /> : <FontAwesomeIcon className='x-button' icon={faXmark} />}</button>
								<p className='loadout-ingredient-amount'>{amount}</p>
								<button className='increase-ingredient' onClick={e => handleIncrementIngredientClick(e, ingredient, 1)}><FontAwesomeIcon icon={faPlus} /></button>
							</li>
						))}
					</ul>
				</div>
			</div>
			<div className='current-loadout-container'>
				<h2 className='loadout-current-label'>Current Loadout {Object.keys(loadout).length > 0 && `(${Object.keys(loadout).length})`}</h2>
				<div className='current-loadout'>
					<h3 className='loadout-total-servings'>Total Servings: {Object.entries(loadout).reduce((acc, [recipeId, amount], index) => acc + (amount * recipes[recipeId].servings), 0)}</h3>
					<h3 className='loadout-total-time'>Total Time Taken: {Object.entries(loadout).reduce((acc, [recipeId, amount], index) => acc + recipes[recipeId].cooktime, 0)} minutes</h3>
					<ul className='loadout-recipes'>
						{Object.entries(loadout).map(([recipeId, amount], index) => (
							<li className='loadout-recipe' key={index}>
								<p className='loadout-recipe-title'>{recipes[recipeId].title}</p>
								<button className='decrease-recipe' onClick={e => handleIncrementRecipeClick(e, recipeId, -1)}>{loadout[recipeId] > 1 ? <FontAwesomeIcon icon={faMinus} /> : <FontAwesomeIcon className='x-button' icon={faXmark} />}</button>
								<p className='loadout-recipe-amount'>{amount}</p>
								<button className='increase-recipe' onClick={e => handleIncrementRecipeClick(e, recipeId, 1)}><FontAwesomeIcon icon={faPlus} /></button>
							</li>
						))}
					</ul>
				</div>
			</div>
			<button type='submit' className='build-loadout'>Build</button>
			<button className='reset-loadout'>Reset</button>
			<div className='query-recipes-main-container'>
				<input type='search' id='recipe-query' name='recipe-query' 
					value={recipeQuery}
					onChange={e => setRecipeQuery(e.target.value)}
				/>
				<div className='query-recipes-filters'>
				</div>
				<div className='query-recipes-container'>
					<table className='query-recipes' border='1'>
						<thead>
							<tr>
								<th>Name</th>
								<th>Servings</th>
								<th>Time</th>
							</tr>
						</thead>
						<tbody>
							{renderQueryRecipes()}
						</tbody>
					</table>
				</div>
			</div>
		</div>
	)
}

export default Loadout
