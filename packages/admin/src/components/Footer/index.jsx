import { useEffect, useRef } from 'react';
import { history, useModel } from 'umi';
import './index.css';
const Footer = () => {
  const { initialState } = useModel('@@initialState');
  const { current } = useRef({ hasInit: false });
  // const version = useMemo(() => {
  //   let v = initialState?.version || '获取中...';
  //   if (history.location.pathname == '/user/login') {
  //     v = '登录后显示';
  //   }
  //   return v;
  // }, [initialState, history]);
  useEffect(() => {
    if (!current.hasInit) {
      current.hasInit = true;
      let v = initialState?.version || '获取中...';
      if (history.location.pathname == '/user/login') {
        v = '登录后显示';
      }
    }
  }, [initialState, history]);
  return null;
  // return (
  //   <>
  //     <div className="footer" style={{ textAlign: 'center', marginTop: 32 }}>
  //       <p>
  //         <span>Powered By </span>
  //         <a className="ua" href="https://vanblog.mereith.com" target="_blank" rel="noreferrer">
  //           VanBlog
  //         </a>
  //       </p>
  //       <p>
  //         <span>版本: </span>
  //         <span> {version}</span>
  //       </p>
  //     </div>
  //   </>
  // );
};

export default Footer;
