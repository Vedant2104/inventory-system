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