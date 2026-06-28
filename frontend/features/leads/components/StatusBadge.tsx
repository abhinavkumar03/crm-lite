export default function StatusBadge({
    status,
}:{
    status:string
}){

    return(

<span
className="rounded bg-gray-100 px-3 py-1 text-sm"
>

{status}

</span>

    );

}