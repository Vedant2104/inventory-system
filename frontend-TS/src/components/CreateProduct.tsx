import type { Product , Category } from "../interfaces/inventory";
import { useState ,useEffect} from "react";
import { useNavigate } from "react-router-dom";
import Loader from "./Loader";
export default function CreateProduct() {

    const [formData, setFormData] = useState<Product>({
        id: "",
        name: "",
        description: "",
        category: { id: "", name: "", description: "" },
        price: "",
        brand: "",
        quantity: "",
      });
    const [isLoading , setIsLoading] = useState<boolean>(false);
    const [error , setError] = useState<string>("");
    const [categories , setCategories] = useState<Category[]>([]);
    const navigate = useNavigate();
    function handleFetchCategories(){
        setIsLoading(true);

        fetch("http://localhost:8080/category", { method: "GET" })
        .then(res=>{
            if(!res.ok){
                throw new Error("Failed to fetch categories");
            }
            return res.json();
        })
        .then(data=>{
            setCategories(data || []);
        }).catch(err=>{
            console.log(err);
            setError(err.message);
        }).finally(()=>{
            setIsLoading(false);
        })
    }

    useEffect(()=>{
        handleFetchCategories();
    },[])

    function handleChange(e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) {
        const {name,value} = e.target;
       
        if(name == 'category'){
            const selected = categories.find((c) => c.id === value);
            setFormData((prev) => prev ? { ...prev, category: selected! } : prev);
        }else if(name == 'price' || name == 'quantity'){
            setFormData((prev) => prev ? { ...prev, [name]: Number(value) } : prev);
        }
        else{
            setFormData((prev) => prev ? { ...prev, [name]: value } : prev);
        }
    }

    function handleSubmit(){
        setIsLoading(true);

        const payload = {
            name : formData?.name,
            description : formData?.description,
            category : formData?.category.id,
            price : formData?.price,
            brand : formData?.brand,
            quantity : formData?.quantity
        }
        console.log(payload);
        fetch("http://localhost:8080/product",{
            method:"POST",
            headers:{
                "Content-Type": "application/json"
            },
            body : JSON.stringify(payload)
        }).then(res=>{
            if(!res.ok){
                throw new Error("Failed to create product");
            }
            return res.json();
        }).then(data=>{
            console.log(data);
            navigate(`/product/${data.id}`);
        }).catch(err=>{
            setError(err.message);
        }).finally(()=>{
            setIsLoading(false);
        })
    }


    function handleBack(){
        navigate("/");
    }

    if(isLoading){
        return <Loader />;
    }
    if(error){
        return <h1>{error}</h1>
    }

  return (
    <div className="product-page">
      <div className="product-card">
      <>
            <div className="card-header">
              <h1 className="product-name">Edit Product</h1>
            </div>

            <div className="form-group">
              <label>Name</label>
              
              <input name="name" value={formData.name} placeholder="Enter product name" onChange={handleChange} />
            </div>

            <div className="form-group">
              <label>Description</label>
              <textarea name="description" value={formData.description} placeholder="Enter product description" onChange={handleChange} />
            </div>

            <div className="form-group">
              <label>Category</label>
              <select name="category" value={formData.category.id} onChange={handleChange}>
                {categories.map((cat) => (
                  <option key={cat.id} value={cat.id}>{cat.name}</option>
                ))}
              </select>
            </div>

            <div className="form-row">
              <div className="form-group">
                <label>Price</label>
                <input name="price" type="number" value={formData.price} placeholder="Enter product price" onChange={handleChange} />
              </div>
              <div className="form-group">
                <label>Brand</label>
                <input name="brand" value={formData.brand} placeholder="Enter brand" onChange={handleChange} />
              </div>
              <div className="form-group">
                <label>Quantity</label>
                <input name="quantity" type="number" value={formData.quantity} placeholder="Enter quantity" onChange={handleChange} />
              </div>
            </div>

            <div className="form-actions">
              <button className="btn-cancel" onClick={handleBack}>Back</button>
              <button className="btn-save" onClick={handleSubmit}>Save</button>
            </div>
          </>
      </div>
    </div>
  )
}
