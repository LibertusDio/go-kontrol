import React, {useEffect} from "react";
import {useParams} from "react-router-dom";
import {useAppDispatch} from "../../app/hooks";
import {getServiceInfo} from "./api";
import {selectService} from "./serviceSlices";
import {useSelector} from "react-redux";


const Service: React.FC = () => {
    const {id} = useParams()
    const dispatch = useAppDispatch();
    // @ts-ignore
    const {statusCode, message} = useSelector(selectService)
    useEffect(() => {
        dispatch(getServiceInfo(id as string))
    })
    return (
        <article className={statusCode === 200 ? "message is-success" : "message is-danger"}>
            <div className="message-header">
                <p>{statusCode}</p>
                <button className="delete" aria-label="delete"></button>
            </div>
            <div className="message-body">{message}</div>
        </article>
    )
}

export {
    Service
}
