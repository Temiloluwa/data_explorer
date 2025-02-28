"use client"
import {redirect} from 'next/navigation'
import React, { useEffect, useState } from 'react';
import {useKindeBrowserClient} from "@kinde-oss/kinde-auth-nextjs";
import { prismadb } from '@/lib/prismadb';
import Dashboard from '@/components/Dashboard';

const Page = () => {
    const { user } = useKindeBrowserClient();
    const [activeUser, setActiveUser ] = useState(null);
    useEffect(() => {
        const checkKindeUserInDb = async () => {
            try {
                // if no kinde user, redirect to callback page  
                console.log(`I am in Dashboards ${user}`);
                if (!user || !user.id) {
                    redirect("/auth-callback?origin=dashboard");
                }

                // if kinde user, get user from database
                const dbUser = await prismadb.user.findUnique({
                    where: {id: user.id}
                });


                // if no db user, redirect to callback page
                if (!dbUser) {
                    redirect("/auth-callback?origin=dashboard");
                }
                
                setActiveUser(() => (setActiveUser(dbUser)));
                
        
            } catch (error) {
                console.log(error)
            }
        
        };
    
        checkKindeUserInDb();
      }, []);

    return <Dashboard user={activeUser}/>

}

export default Page;
