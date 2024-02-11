
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import SignIn from "./pages/sign-in/SignIn"
import SignUp from "./pages/sign-up/SignUp"
export default function App() {
  
  if (!localStorage.getItem('userRole'))
    localStorage.setItem('userRole', 'guest');
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<SignUp />} />
        <Route path="/sign-in" element={<SignIn />} />
        <Route path="/sign-up" element={<SignUp />} />
      </Routes>
    </BrowserRouter>
  )
}
