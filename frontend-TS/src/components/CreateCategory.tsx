import type { Category } from "../interfaces/inventory";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import Loader from "./Loader";
export default function CreateCategory() {

    const [formData, setFormData] = useState<Category>({
        id: "",
        name: "",
        description: "",
      });
    const [isLoading , setIsLoading] = useState<boolean>(false);
    const [error , setError] = useState<string>("");
    const navigate = useNavigate();
    

    function handleChange(e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) {
        const { name, value } = e.target;
        setFormData((prev) => prev ? { ...prev, [name]: value } : prev);
    }

    function handleSubmit(){
        setIsLoading(true);

        const payload = {
            name : formData?.name,
            description : formData?.description,
        }
        console.log(payload);
        fetch("http://localhost:8080/category",{
            method:"POST",
            headers:{
                "Content-Type": "application/json"
            },
            body : JSON.stringify(payload)
        }).then(res=>{
            if(!res.ok){
                throw new Error("Failed to create category");
            }
            return res.json();
        }).then(data=>{
            console.log(data);
            navigate(`/category/${data.id}`);
        }).catch(err=>{
            setError(err.message);
        }).finally(()=>{
            setIsLoading(false);
        })
    }


    function handleBack(){
        navigate("/category");
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
              
              <input name="name" value={formData.name} placeholder="Enter category name" onChange={handleChange} />
            </div>

            <div className="form-group">
              <label>Description</label>
              <textarea name="description" value={formData.description} placeholder="Enter category description" onChange={handleChange} />
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
