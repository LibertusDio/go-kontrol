import React, {useEffect, useState} from 'react';
import {useForm, SubmitHandler} from 'react-hook-form';
import {useSelector, useDispatch} from 'react-redux';
import {AuthModel, AuthRequest, AuthRequestInterface} from "./AuthModel";
import {clearState, login, userSelector} from "./login-slice";
import {useNavigate} from "react-router-dom";

export function Login() {
    const navigate = useNavigate()
    const dispatch = useDispatch();
    const {isLoggedIn} = useSelector(userSelector)
    const {register, handleSubmit, watch, formState: {errors}} = useForm<AuthRequest>();
    const onSubmit: SubmitHandler<AuthRequest> = (data: AuthRequest) => {
        // @ts-ignore
        dispatch(login(data))
    };
    useEffect(() => {
        if (isLoggedIn === true) {
            navigate("/")
        }
    }, [isLoggedIn])
    return (
        <section className="section">
            <form onSubmit={handleSubmit(onSubmit)}>
                <div className="field">
                    <label className="label">Service ID</label>
                    <div className="control">
                        <select className="input" {...register("service_id")}>
                            <option value="88503398-db0f-11ec-9d64-0242ac120002">idt</option>
                            <option value="88503398-db0f-11ec-9d64-0242ac120003">adt</option>
                            <option value="88503398-db0f-11ec-9d64-0242ac120004">hrd</option>
                        </select>
                    </div>
                </div>
                <div className="field">
                    <label className="label">User name</label>
                    <div className="control">
                        <input className="input" type="text" {...register("user_name")}/>
                    </div>
                </div>
                <div className="field">
                    <label className="label">Password</label>
                    <div className="control">
                        <input className="input" type="password" {...register("password")}/>
                    </div>
                </div>
                <div className="field is-grouped">
                    <div className="control">
                        <button className="button is-link">Submit</button>
                    </div>
                </div>
            </form>
        </section>
    )
}
