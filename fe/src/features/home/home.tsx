import React, {useEffect, useState} from 'react';
import {useNavigate} from "react-router-dom"
import {useDispatch, useSelector} from "react-redux";
import {clearState, userSelector} from "../login/login-slice";
import {serviceInfo, serviceSelector} from "./serviceSlice";
import {useAppSelector} from "../../app/hooks";

export function Home() {
    const navigate = useNavigate()
    const dispatch = useDispatch()
    const {isLoggedIn} = useAppSelector(userSelector)
    const {statusCode, effectedAt, message} = useAppSelector(serviceSelector)

    useEffect(() => {
        if (isLoggedIn === false) {
            // dispatch(clearState());
            navigate('/login');
        }
    }, [isLoggedIn]);
    const handleAction = (serviceId: string) => {
        // @ts-ignore
        dispatch(serviceInfo(serviceId))
    }
    useEffect(() => {
        switch (statusCode) {
            case 401:
                dispatch(clearState());
                //navigate("/login")
                break
            case 403:
                alert(`Status code: ${statusCode} - Message: ${message}`)
                break
            case 200:
                alert(`Status code: ${statusCode} - Message: ${message}`)
                break
        }
    }, [statusCode, effectedAt, message])
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
