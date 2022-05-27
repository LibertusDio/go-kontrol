import React, {useEffect, useState} from 'react';
import {useNavigate} from "react-router-dom"
import {useDispatch, useSelector} from "react-redux";
import {clearState, userSelector} from "../login/login-slice";
import {serviceInfo} from "./serviceSlice";

export function Home() {
    const navigate = useNavigate()
    const dispatch = useDispatch()
    const {isLoggedIn} = useSelector(userSelector)
    console.log(process.env.DOMAIN_API_GATEWAY);
    
    useEffect(() => {
        if (isLoggedIn === false) {
            dispatch(clearState());
            navigate('/login');
        }
    }, [isLoggedIn]);
    const handleAction = (serviceId: string) => {
        // @ts-ignore
        dispatch(serviceInfo(serviceId))
    }
    return (
        <section className="section">
            <div className="columns">
                <div className="buttons column is-three-fifths is-offset-one-fifth">
                    <button className="button is-info" onClick={() => handleAction('idt')}>idt</button>
                    <button className="button is-warning" onClick={() => handleAction('adt')}>adt</button>
                    <button className="button is-danger" onClick={() => handleAction('hrd')}>hrd</button>
                </div>
            </div>
        </section>
    )
}
