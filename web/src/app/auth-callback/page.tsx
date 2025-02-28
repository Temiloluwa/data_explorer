'use client'
import { useRouter, useSearchParams } from 'next/navigation'
import React, { useEffect } from 'react';
import axios from "axios";
import { Loader2 } from 'lucide-react';

/* The purpose of this page is following:

1. Call the auth-callback endpoint to check if the user is synced to db
2. If the user is synced, redirect to the origin page
3. If the user is not synced, redirect to the login page

*/
const Page = () => {
    const router = useRouter()
    const searchParams = useSearchParams()
    const origin = searchParams.get("origin");
    useEffect(() => {
        const syncUserWithDb = async () => {
            try {
                const {data} = await axios.get<{ message?: string; success: boolean; }>("/api/auth/auth-callback");
                console.log(`I am in client page auth-callback ${data}`)
        
                if (data.success) {
                    // user is synced to db, go back to the dashboard or the origin you are coming from 
                    router.push( origin ? `/${origin}`: "/dashboard")
                } else {
                    // user is not synced to db
                    console.log("User is not synced to db", data)
                    router.push("/")
                }
        
            } catch (error) {
                console.log(error)
            }
        
        };
    
        syncUserWithDb();
      }, []);


    return (
        <div className='w-full m-24 flex justify-center'>
            <div className='flex flex-col items-center gap-2'>
                <Loader2 className='h-8 w-8 animate-spin text-zinc-800'/>
                <h3 className='font-semibold text-xl'>Setting up your account ...</h3>
                <p> You will be redirected automatically</p>
            </div>
        </div>
    )
        
}

export default Page;