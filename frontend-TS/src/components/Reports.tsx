import LowStockReport from "./Reports/LowStockReport";
import ProductCountByCategoryReport from "./Reports/ProductCountByCategoryReport";
import { useState } from "react";
import "../styles/Reports.css";
import PriceSegmentationReport from "./Reports/PriceSegmentationReport";

const WINDOWS: string[] = [
  "Low Stock Report",
  "Product Count By Category Report",
  "Price Segmentation Report",
];

export default function Reports() {
  const [active, setActive] = useState<string>("Low Stock Report");

  return (
    <div className="reports-container">
      <h2 className="reports-title">Reports Dashboard</h2>

      <div className="reports-tabs">
        {WINDOWS.map((item) => (
          <button
            key={item}
            className={`tab-btn ${active === item ? "active" : ""}`}
            onClick={() => setActive(item)}
          >
            {item}
          </button>
        ))}
      </div>

      <div className="reports-content">
        {active === "Low Stock Report" && <LowStockReport />}
        {active === "Product Count By Category Report" && (
          <ProductCountByCategoryReport />
        )}
        {active === "Price Segmentation Report" && <PriceSegmentationReport/>}
      </div>
    </div>
  );
}