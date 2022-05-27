import React, {useEffect, useState} from 'react';
import {useNavigate} from "react-router-dom";
import {useDispatch, useSelector} from "react-redux";
import {clearState, userSelector} from "../login/login-slice";
export function Header() {
    const {isLoggedIn} = useSelector(userSelector)
    console.log(isLoggedIn)
    return (
        <nav className="navbar" role="navigation" aria-label="main navigation">
            <div className="navbar-brand">

            </div>

            <div id="navbarBasicExample" className="navbar-menu">
                <div className="navbar-start">

                </div>

                <div className="navbar-end">
                    <div className="navbar-item">
                        <div className="buttons">
                            <a className="button is-primary">
                                <strong>Sign up</strong>
                            </a>
                            <a className="button is-light">
                                Log in
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        </nav>
    )
}
