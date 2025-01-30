import { createBrowserRouter, createRoutesFromElements, Route } from 'react-router-dom';
import { ROUTES } from "./routes"

const router = createBrowserRouter(
	createRoutesFromElements(
		ROUTES.map((route) => (
			<Route key={route.path} path={route.path} element={route.element} />
		))
	)
)
		
export default router
