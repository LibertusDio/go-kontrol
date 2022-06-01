import React, {useEffect} from "react";
import {useLocation, useNavigate} from "react-router-dom";
import {useForm, SubmitHandler} from 'react-hook-form';
import {useAppDispatch} from "../../app/hooks";
import {AuthRequest} from "./authModel";
import {login} from "./api";
import {useSelector} from "react-redux";
import {userSelector} from "./loginSlices";


const Login: React.FC = (props: any) => {
    const dispatch = useAppDispatch();
    const {handleSubmit, register} = useForm<AuthRequest>();
    let location = useLocation()
    let navigate = useNavigate()
    // @ts-ignore
    const { from } = location.state || { from: { pathname: "/" } };
    const {isLoggedIn} = useSelector(userSelector)
    const onSubmit: SubmitHandler<AuthRequest> = (data: AuthRequest) => {
        // @ts-ignore
        dispatch(login(data))
    };
    useEffect(() => {
        if (isLoggedIn) {
            navigate(from as string)
        }
    },[isLoggedIn])

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

export
{
    Login
}
