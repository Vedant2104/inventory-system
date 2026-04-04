import './App.css'
import Navbar from './components/Navbar'
import ProductList from './components/ProductList'
import { BrowserRouter , Route , Routes } from 'react-router-dom'
import ProductPage from './components/ProductPage'
import CategoryList from './components/CategoryList'
import CategoryPage from './components/CategoryPage'
import CreateProduct from './components/CreateProduct'
import CreateCategory from './components/CreateCategory'

function App() {

  return (
    <>
      <BrowserRouter>
        <Navbar/> 
        <Routes>
          <Route path="/" element={<ProductList />}/>
          <Route path="/category" element={<CategoryList/>}/>
          <Route path="/product/new" element={<CreateProduct/>}/>
          <Route path="/product/:id" element={<ProductPage/>}/>
          <Route path="/category/new" element={<CreateCategory/>}/>
          <Route path="/category/:id" element={<CategoryPage/>}/>
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
