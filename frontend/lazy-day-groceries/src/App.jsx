import './App.css'
import { RouterProvider } from 'react-router-dom';
import NavigationBar from './common/NavigationBar.jsx'
import router from './routes/router.jsx'

function App() {
	return (
		<>
			<NavigationBar />
			<RouterProvider router={router} />
		</>
	)
  /*return (
    <>
	    <NavigationBar />
	    <HomeContent />
    </>
  )
  */
}

export default App
