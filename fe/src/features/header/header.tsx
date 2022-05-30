import React, {useEffect, useState} from 'react';
import {useNavigate} from "react-router-dom";
import {useDispatch, useSelector} from "react-redux";
import {clearState, userSelector} from "../login/login-slice";

export function Header() {
    const {isLoggedIn} = useSelector(userSelector)
    return (
        <nav className="navbar" role="navigation" aria-label="main navigation">
            <div className="navbar-brand">

            </div>

            <div id="navbarBasicExample" className="navbar-menu">
                <div className="navbar-start">

                </div>

                <div className="navbar-end">
                    {isLoggedIn ? <Logout/> : null}
                </div>
            </div>
        </nav>
    )
}

const Logout = () => {
    const dispatch = useDispatch()
    const logout = () => {
        dispatch(clearState())
    }
    return (
        <div className="navbar-item">
            <div className="buttons">
                <a className="button is-primary" onClick={() => logout()}>
                    <strong>Logout</strong>
                </a>
            </div>
        </div>
    )
}
