
export const CSVDownloader =<T extends object>(data : T[] , filename : string) => {
    if(data.length === 0){
        return Error("No data to download");
    }

    const headers = Object.keys(data[0]).join(',');

    const rows = data.map(obj=>
        Object.values(obj).map(value => `"${String(value).replace(/"/g , '"')}"`)
        .join(',')
    );
    console.log(rows);
    const csvContent = [headers,...rows].join('\n');
    const blob = new Blob([csvContent] , {type : "text/csv;charset=utf-8;"})

    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);
    link.setAttribute('href' , url)
    link.setAttribute('download' , `${filename}.csv`)
    link.style.visibility = "hidden";
    document.body.appendChild(link);
    link.click()
    document.body.removeChild(link)
}