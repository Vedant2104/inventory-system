export interface Category {
    id : string;
    name : string;
    description : string;
}

export interface Product{
    id: string;
    name : string;
    description: string
    category : Category;
    price :number|"" ;
    brand : string;
    quantity : number|"";
}

export interface ProductCountByCategory{
    category_id : string;
    category : string;
    count : number;
}

export interface LowStockProducts{
    id : string;
    product : string;
    brand : string;
    quantity : number;
}

export interface PriceSegmentation{
    category_id : string;
    category : string;
    budget : number;
    midRange : number;
    premium : number;
}