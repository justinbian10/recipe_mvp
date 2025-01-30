import Home from "../pages/Home/Home.jsx"
import RecipeView from "../pages/Recipe/RecipeView.jsx"
import RecipeEdit from "../pages/Recipe/RecipeEdit.jsx"
import Loadout from "../pages/Loadout/Loadout.jsx"

export const ROUTES = [
	{
		path: "/",
		element: <Home />,
	},
	{
		path: "/recipe/:recipeId",
		element: <RecipeView />,
	},
	{
		path: "/recipe/:recipeId/edit",
		element: <RecipeEdit />,
	},
	{
		path: "/loadout",
		element: <Loadout />,
	}
]

