import { useState } from "react"
import type { LowStockProducts } from "../../interfaces/inventory"
import Loader from "../Loader";
import { Link } from "react-router-dom";
import "../../styles/lowStockReport.css"
import { CSVDownloader } from "../../utils/CSVDowloader";
export default function LowStockReport() {
    const [products, setProducts] = useState<LowStockProducts[]>([])
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [threshold, setThreshold] = useState<string>("5");
    const [error, setError] = useState<string>("");
    const [hasGenerated, setHasGenerated] = useState<boolean>(false);  // track if report was generated

    function handleFetchProducts() {
        if (!threshold || Number(threshold) < 0) {   
            setError("Please enter a valid threshold");
            return;
        }

        setIsLoading(true);
        setError("");           
        setHasGenerated(false);

        fetch(`http://localhost:8080/product/report/lowstock?threshold=${threshold}`, {
            method: "GET",
        })
        .then(resp => {
            if (!resp.ok) throw new Error("Failed to fetch products");
            return resp.json();
        })
        .then(data => {
            setProducts(data || []);
            setHasGenerated(true);
        })
        .catch(err => {
            setError(err.message);
        })
        .finally(() => {
            setIsLoading(false);
        });
    }

    function handleKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
        if (e.key === "Enter") handleFetchProducts();   
    }

    if (isLoading) return <Loader />;

    return (
        <div className="lowstock-page">
            <button className={`btn-download ${!hasGenerated ? "disabled" : ""}`} disabled={!hasGenerated} onClick={() => CSVDownloader(products , `products-below-threshold-${threshold}`)}>Download</button>
            <div className="lowstock-container">

                <div className="lowstock-header">
                    <h2 className="lowstock-title">Low Stock Report</h2>
                    <p className="lowstock-subtitle">Products below the quantity threshold</p>
                </div>

                <div className="lowstock-controls">
                    <div className="input-group">
                        <label htmlFor="threshold">Threshold</label>
                        <input
                            id="threshold"
                            type="number"
                            value={threshold}
                            onChange={(e) => setThreshold(e.target.value)}
                            onKeyDown={handleKeyDown}
                            min="0"
                        />
                    </div>
                    <button className="btn-generate" onClick={handleFetchProducts}>
                        Generate Report
                    </button>
                </div>

                {error && <p className="lowstock-error">{error}</p>}

                {/* only show "no products" after a report has been generated */}
                {hasGenerated && products.length === 0 && (
                    <p className="lowstock-empty">No products found below threshold.</p>
                )}

                {products.length > 0 && (
                    <>
                        <p className="lowstock-count">{products.length} product(s) found</p>
                        <table className="lowstock-table">
                            <thead>
                                <tr>
                                    <th>Product</th>
                                    <th>Brand</th>
                                    <th>Quantity</th>
                                </tr>
                            </thead>
                            <tbody>
                                {products.map(item => (
                                    <tr key={item.id} className={item.quantity === 0 ? "out-of-stock" : ""}>
                                        <td>
                                            <Link to={`/product/${item.id}`}>{item.product}</Link>
                                        </td>
                                        <td>{item.brand}</td>
                                        <td>
                                            <span className={`quantity-badge ${item.quantity === 0 ? "badge-empty" : "badge-low"}`}>
                                                {item.quantity}
                                            </span>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </>
                )}

            </div>
        </div>
    );
}