import "./App.css";
import {notification} from "antd";
import {Navigate, Route, Routes, useLocation} from "react-router-dom";
import {HomesteadHeader} from "./containers/header/Header";
import {useEffect} from "react";
import {useAppSelector} from "./app/hooks";
import {Service} from "./features/service/service";
import {Login} from "./features/login/login";
import Home from "./features/home/home"

function App() {
    const error = useAppSelector((state) => state.error);
    useEffect(() => {
        if (error.message !== "") {
            openNotification(error.message);
        }
    });
    const openNotification = (message: string) => {
        notification.error({
            message: `Error`,
            description: message,
        });
    };
    return (
        <div className="container">
            <HomesteadHeader/>
            <Routes>
                <Route path=":id/service" element={
                    <RequireAuth>
                        <Service/>
                    </RequireAuth>
                }/>
                <Route path="login" element={<Login/>}/>
                <Route path="/" element={<Home/>}/>
            </Routes>
        </div>
    );
}

function RequireAuth({children}: { children: JSX.Element }) {
    const isLogin = useAppSelector((state) => state.auth.isLoggedIn);
    const location = useLocation();

    if (!isLogin) {
        // Redirect them to the /login page, but save the current location they were
        // trying to go to when they were redirected. This allows us to send them
        // along to that page after they login, which is a nicer user experience
        // than dropping them off on the home page.
        return <Navigate to="/login" state={{from: location}} replace/>;
    }

    return children;
}

export default App;
