(this.webpackJsonprtsp2rtmpweb=this.webpackJsonprtsp2rtmpweb||[]).push([[0],{165:function(e,a,t){"use strict";t.r(a);var r=t(0),n=t.n(r),c=t(10),o=t.n(c),i=t(49),l=t(13),d=t(203),s=t(167),j=t(217),b=t(221),m=t(220),u=t(216),h=t(218),p=t(237),x=t(219),O=t(110),g=t.n(O),f=t(109),v=t.n(f),w=t(215),k=t(74),C=t.n(k),y=t(209),S=t(239),R=t(207),P=t(208),I=t(37),N=t(210),B=t(48),T=t.n(B),L=t(205),E=t(44),M=t.n(E);let D;D=window.location.origin;var V={serverURL:D},Y=t(32),U=t.n(Y);U.a.defaults.withCredentials=!1,U.a.defaults.timeout=1e4,U.a.interceptors.request.use((e=>{let a=window.localStorage.getItem("token");return e.headers.token=a,e}),(e=>(console.error(e),Promise.reject(e)))),U.a.interceptors.response.use((function(e){return e.data&&e.data.errcode&&401===parseInt(e.data.errcode)&&(window.location.hash="#/login"),e}),(function(e){return U.a.isCancel()||(e.response?(console.log(401===e.response.status),401===e.response.status?window.localStorage.getItem("token")&&window.localStorage.getItem("tokenExpired")&&"false"!==window.localStorage.getItem("tokenExpired")||(window.localStorage.setItem("tokenExpired","true"),window.location.hash="#/login"):500===e.response.status&&alert("server exception !")):e&&"error: timeout"===String(e).toLowerCase().substring(0,14)?alert("server timeout !"):alert("server error !")),Promise.reject(e)}));let H=V.serverURL;const F=(e,a)=>U.a.post(`${H}${e}`,a).then((e=>e.data)),A=(e,a)=>U.a.get(`${H}${e}`,{params:a}).then((e=>e.data));var W=`${V.serverURL}`,z=e=>F("/system/login",e),$=e=>A("/camera/list",e),q=e=>A("/camera/detail",e),J=e=>F("/camera/edit",e),_=e=>F(`/camera/delete/${e.id}`,e),K=e=>F("/camera/enabled",e),G=e=>F("/camera/savevideochange",e),Q=e=>F("/camera/livechange",e),X=e=>F("/camera/rtmppushchange",e),Z=e=>F("/camera/playauthcodereset",e),ee=e=>A("/camerashare/list",e),ae=e=>F("/camerashare/edit",e),te=e=>F(`/camerashare/delete/${e.id}`,e),re=e=>F("/camerashare/enabled",e),ne=t(2);const ce=Object(d.a)((e=>({appBar:{position:"relative",backgroundColor:"#eebbaa"},title:{marginLeft:e.spacing(2),flex:1},videoContainer:{width:"90%",margin:"0 auto"}}))),oe=n.a.forwardRef((function(e,a){return Object(ne.jsx)(L.a,{direction:"up",ref:a,...e})}));function ie(e){const a=ce(),[t,c]=n.a.useState(!1),[o,i]=n.a.useState(!0),[l,d]=n.a.useState(e.playParam||{playMethod:"",cameraCode:"",authCode:""});var s=null,j=0;Object(r.useImperativeHandle)(e.onRef,(()=>({handleClickOpen:b})));const b=()=>{c(!0)},m=()=>{c(!1)},u=e=>{var a={type:"flv"};let t=W+"/live/"+l.playMethod+"/"+l.cameraCode+"/"+l.authCode+".flv";a.url=t,a.hasAudio=o,a.isLive=!0,console.log("MediaDataSource",a),h(a)},h=e=>{var a=document.getElementsByClassName("centeredVideo")[0];"undefined"!==typeof s&&null!=s&&(s.pause(),s.unload(),s.detachMediaElement(),s.destroy(),s=null),(s=M.a.createPlayer(e,{enableWorker:!1,lazyLoadMaxDuration:180,seekType:"range"})).on(M.a.Events.ERROR,((e,a,t)=>{console.log("errorType:",e),console.log("errorDetail:",a),console.log("errorInfo:",t),s&&(s.pause(),s.unload(),s.detachMediaElement(),s.destroy(),s=null,window.setTimeout(u,500))})),s.on(M.a.Events.STATISTICS_INFO,(function(e){0!=j?j!=e.decodedFrames?j=e.decodedFrames:(console.log("decodedFrames:",e.decodedFrames),j=0,s&&(s.pause(),s.unload(),s.detachMediaElement(),s.destroy(),s=null,window.setTimeout(u,500))):j=e.decodedFrames})),s.attachMediaElement(a),s.load(),s.play()};return Object(ne.jsx)("div",{children:Object(ne.jsxs)(S.a,{fullScreen:!0,open:t,onClose:m,TransitionComponent:oe,children:[Object(ne.jsx)(R.a,{className:a.appBar,children:Object(ne.jsxs)(P.a,{children:[Object(ne.jsx)(y.a,{variant:"contained",onClick:u,children:"play"}),Object(ne.jsxs)(I.a,{variant:"h6",className:a.title,children:["hasAudio",Object(ne.jsx)(N.a,{checked:o,id:"Audio",color:"primary",name:"hasAudio",onChange:e=>{i(e.target.checked)},inputProps:{"aria-label":"primary checkbox"}})]}),Object(ne.jsx)(y.a,{autoFocus:!0,color:"inherit",onClick:m,children:Object(ne.jsx)(T.a,{})})]})}),Object(ne.jsx)("div",{className:a.videoContainer,children:Object(ne.jsx)("div",{children:Object(ne.jsx)("video",{name:"videoElement",className:"centeredVideo",controls:!0,allow:"autoPlay",width:"100%",children:"Your browser is too old which doesn't support HTML5 video."})})})]})})}var le=t(236),de=t(234),se=t(211);const je=Object(d.a)((e=>({appBar:{position:"relative"},title:{marginLeft:e.spacing(2),flex:1},formClass:{"& > *":{margin:e.spacing(1),width:"25ch"}},formDiv:{margin:"0 auto"}}))),be=n.a.forwardRef((function(e,a){return Object(ne.jsx)(L.a,{direction:"up",ref:a,...e})}));function me(e){const a=je(),[t,c]=n.a.useState(!1),[o,i]=n.a.useState(1===e.row.enabled),[l,d]=n.a.useState(1===e.row.saveVideo),[s,j]=n.a.useState(1===e.row.live),[b,m]=n.a.useState(1===e.row.live),[u,h]=n.a.useState({id:e.row.id,code:e.row.code,rtspURL:e.row.rtspURL,rtmpURL:e.row.rtmpURL,playAuthCode:e.row.playAuthCode,onlineStatus:e.row.onlineStatus,enabled:e.row.enabled,saveVideo:e.row.saveVideo,live:e.row.live,rtmpPushStatus:e.row.rtmpPushStatus}),[p,x]=n.a.useState(!1),[O,g]=n.a.useState("");Object(r.useImperativeHandle)(e.onRef,(()=>({handleClickOpen:f})));const f=()=>{c(!0)},v=()=>{c(!1)},w=a=>{J(u).then((a=>{if(1===a.code)return c(!1),void(e.callBack&&e.callBack());g(a.msg),x(!0),window.setTimeout((function(){x(!1)}),5e3)}))},k=e=>{u[e.target.id]=e.target.value},C=e=>{console.log(e.target.checked),u[e.target.id]=e.target.checked?1:0,"enabled"===e.target.id&&i(e.target.checked),"saveVideo"===e.target.id&&d(e.target.checked),"live"===e.target.id&&j(e.target.checked),"rtmpPushStatus"===e.target.id&&m(e.target.checked)};return Object(ne.jsx)("div",{children:Object(ne.jsxs)(S.a,{fullScreen:!0,open:t,onClose:v,TransitionComponent:be,children:[Object(ne.jsx)(R.a,{className:a.appBar,children:Object(ne.jsxs)(P.a,{children:[Object(ne.jsx)(y.a,{variant:"contained",onClick:w,children:"\u4fdd\u5b58"}),Object(ne.jsx)("span",{children:"\xa0\xa0"}),"edit"===e.type?Object(ne.jsx)(y.a,{variant:"contained",onClick:a=>{_(u).then((a=>{if(1===a.code)return c(!1),void(e.callBack&&e.callBack());g(a.msg),x(!0),window.setTimeout((function(){x(!1)}),5e3)}))},children:"\u5220\u9664"}):"",Object(ne.jsx)(I.a,{variant:"h6",className:a.title}),Object(ne.jsx)(y.a,{autoFocus:!0,color:"inherit",onClick:v,children:Object(ne.jsx)(T.a,{})})]})}),p?Object(ne.jsxs)(de.a,{severity:"error",children:[Object(ne.jsx)(se.a,{children:"Error"}),O," ",Object(ne.jsx)("strong",{children:"check it out!"})]}):"",Object(ne.jsx)("div",{className:a.formDiv,children:Object(ne.jsxs)("form",{className:a.formClass,noValidate:!0,autoComplete:"off",onSubmit:w,children:["edit"===e.type?Object(ne.jsx)("div",{children:Object(ne.jsx)(le.a,{id:"id",label:"id",InputProps:{readOnly:!0},defaultValue:u.id})}):"",Object(ne.jsx)("div",{children:Object(ne.jsx)(le.a,{id:"code",label:"\u6444\u50cf\u5934\u7f16\u53f7",defaultValue:u.code,onChange:k})}),Object(ne.jsx)("div",{children:Object(ne.jsx)(le.a,{id:"rtspURL",label:"rtspURL",defaultValue:u.rtspURL,onChange:k})}),Object(ne.jsx)("div",{children:Object(ne.jsx)(le.a,{id:"rtmpURL",label:"rtmpURL",defaultValue:u.rtmpURL,onChange:k})}),"edit"===e.type?"":Object(ne.jsxs)("div",{children:["\u542f\u7528\u72b6\u6001\uff1a",Object(ne.jsx)(N.a,{disabled:!0,checked:o,id:"enabled",onChange:C,color:"primary",name:"enabled",inputProps:{"aria-label":"primary checkbox"}})]}),Object(ne.jsxs)("div",{children:["\u5f55\u50cf\u72b6\u6001\uff1a",Object(ne.jsx)(N.a,{disabled:!0,checked:l,id:"saveVideo",onChange:C,color:"primary",name:"saveVideo",inputProps:{"aria-label":"primary checkbox"}})]}),Object(ne.jsxs)("div",{children:["\u76f4\u64ad\u72b6\u6001\uff1a",Object(ne.jsx)(N.a,{disabled:!0,checked:s,id:"live",onChange:C,color:"primary",name:"live",inputProps:{"aria-label":"primary checkbox"}})]}),Object(ne.jsxs)("div",{children:["Rtmp\u63a8\u9001\u72b6\u6001\uff1a",Object(ne.jsx)(N.a,{disabled:!0,checked:b,id:"rtmpPushStatus",onChange:C,color:"primary",name:"rtmpPushStatus",inputProps:{"aria-label":"primary checkbox"}})]})]})})]})})}const ue=Object(d.a)((()=>({root:{position:"relative"},dropdown:{position:"absolute",top:28,right:0,left:0,zIndex:999,border:"0px solid",padding:1,backgroundColor:"#f8f8f8"}})));function he(e){const a=ue(),[t,r]=n.a.useState(!1),[c,o]=n.a.useState(e.row),l=()=>{c.enabled=1===c.enabled?0:1,K(c).then((a=>{if(1===a.code)return r(!1),void(e.callBack&&e.callBack())}))},d=()=>{c.saveVideo=1===c.saveVideo?0:1,G(c).then((a=>{if(1===a.code)return r(!1),void(e.callBack&&e.callBack())}))},s=()=>{c.rtmpPushStatus=1===c.rtmpPushStatus?0:1,X(c).then((a=>{if(1===a.code)return r(!1),void(e.callBack&&e.callBack())}))},j=()=>{c.live=1===c.live?0:1,Q(c).then((a=>{if(1===a.code)return r(!1),void(e.callBack&&e.callBack())}))};let b=n.a.createRef();let m=n.a.createRef();return Object(ne.jsx)(w.a,{onClickAway:()=>{r(!1)},children:Object(ne.jsxs)("div",{className:a.root,children:[Object(ne.jsx)("button",{type:"button",onClick:()=>{r((e=>!e))},children:Object(ne.jsx)(C.a,{})}),t?Object(ne.jsx)("div",{className:a.dropdown,children:Object(ne.jsxs)("div",{style:{display:"flex",flexDirection:"column",gap:"5px"},children:[Object(ne.jsx)("button",{onClick:function(){r(!1),m.current.handleClickOpen()},children:"\u7f16\u8f91"}),1===e.row.enabled?Object(ne.jsx)("button",{onClick:l,children:"\u7981\u7528"}):Object(ne.jsx)("button",{onClick:l,children:"\u542f\u7528"}),Object(ne.jsx)("button",{onClick:function(){r(!1),b.current.handleClickOpen()},children:"\u64ad\u653e"}),Object(ne.jsx)("button",{children:Object(ne.jsx)(i.b,{exact:!0,to:{pathname:"/camerashare",search:"?cameraId="+e.row.id,state:{row:e.row}},children:"\u5206\u4eab"})}),Object(ne.jsx)("button",{onClick:()=>{Z(c).then((a=>{if(1===a.code)return r(!1),void(e.callBack&&e.callBack())}))},children:"\u91cd\u7f6e\u64ad\u653e\u6743\u9650\u7801"}),1===e.row.saveVideo?Object(ne.jsx)("button",{onClick:d,children:"\u505c\u6b62\u5f55\u50cf"}):Object(ne.jsx)("button",{onClick:d,children:"\u5f00\u542f\u5f55\u50cf"}),1===e.row.live?Object(ne.jsx)("button",{onClick:j,children:"\u505c\u6b62\u76f4\u64ad"}):Object(ne.jsx)("button",{onClick:j,children:"\u5f00\u542f\u76f4\u64ad"}),1===e.row.rtmpPushStatus?Object(ne.jsx)("button",{onClick:s,children:"\u505c\u6b62Rtmp\u63a8\u9001"}):Object(ne.jsx)("button",{onClick:s,children:"\u5f00\u542fRtmp\u63a8\u9001"})]})}):null,Object(ne.jsx)(ie,{playParam:{playMethod:"permanent",cameraCode:e.row.code,authCode:e.row.playAuthCode},onRef:b}),Object(ne.jsx)(me,{row:e.row,type:"edit",callBack:e.callBack,onRef:m})]})})}const pe=[{id:"id",label:"id",minWidth:20},{id:"code",label:"\u6444\u50cf\u5934\u7f16\u53f7",minWidth:20},{id:"rtspURL",label:"rtspURL"},{id:"rtmpURL",label:"rtmpURL"},{id:"playAuthCode",label:"\u64ad\u653e\u6743\u9650\u7801"},{id:"onlineStatus",label:"\u5728\u7ebf\u72b6\u6001",format:e=>e&&1===e?Object(ne.jsx)(v.a,{}):Object(ne.jsx)(g.a,{})},{id:"enabled",label:"\u542f\u7528\u72b6\u6001",format:e=>Object(ne.jsx)(N.a,{checked:1===e,id:"enabled",color:"primary",name:"enabled",inputProps:{"aria-label":"primary checkbox"}})},{id:"live",label:"\u76f4\u64ad\u72b6\u6001",format:e=>Object(ne.jsx)(N.a,{checked:1===e,id:"live",color:"primary",name:"live",inputProps:{"aria-label":"primary checkbox"}})},{id:"saveVideo",label:"\u5f55\u50cf\u72b6\u6001",format:e=>Object(ne.jsx)(N.a,{checked:1===e,id:"saveVideo",color:"primary",name:"saveVideo",inputProps:{"aria-label":"primary checkbox"}})},{id:"rtmpPushStatus",label:"Rtmp\u63a8\u9001\u72b6\u6001\uff1a",format:e=>Object(ne.jsx)(N.a,{checked:1===e,id:"rtmpPushStatus",color:"primary",name:"rtmpPushStatus",inputProps:{"aria-label":"primary checkbox"}})},{id:"action",label:"\u64cd\u4f5c",minWidth:150,format:(e,a,t)=>Object(ne.jsx)(he,{row:a,callBack:t})}],xe=Object(d.a)((e=>({root:{width:"100%",position:"relative"},container:{minHeight:400},dropdown:{position:"absolute",top:28,right:0,left:0,zIndex:0,border:"0px solid",padding:1,backgroundColor:"#f8f8f8"},appBar:{position:"relative"},title:{marginLeft:e.spacing(2),flex:1}})));function Oe(){const e=xe(),[a,t]=n.a.useState(0),[r,c]=n.a.useState(10);var o,i;[o,i]=n.a.useState([]);const l=()=>{$().then((e=>{1===e.code&&(o.splice(0),o.push(...e.data.page),i([]),i(o))}))};let d=n.a.createRef();return n.a.useEffect(l,[]),Object(ne.jsxs)(s.a,{className:e.root,children:[Object(ne.jsx)(R.a,{className:e.appBar,children:Object(ne.jsxs)(P.a,{children:[Object(ne.jsx)(y.a,{variant:"contained",onClick:function(){d.current.handleClickOpen()},children:"\u521b\u5efa"}),Object(ne.jsx)(I.a,{variant:"h6",className:e.title,children:"\u6444\u50cf\u5934\u5217\u8868"})]})}),Object(ne.jsx)(me,{row:{id:"",code:"",rtspURL:"rtsp://192.168.0.10:554/1",rtmpURL:"rtmp://127.0.0.1:1935/code/authcode",playAuthCode:"",onlineStatus:0,enabled:1,saveVideo:0,live:1,rtmpPushStatus:1},type:"add",callBack:l,onRef:d}),Object(ne.jsx)(u.a,{className:e.container,children:Object(ne.jsxs)(j.a,{stickyHeader:!0,"aria-label":"sticky table",children:[Object(ne.jsx)(h.a,{children:Object(ne.jsx)(x.a,{children:pe.map((e=>Object(ne.jsx)(m.a,{align:e.align,style:{minWidth:e.minWidth},children:e.label},e.id)))})}),Object(ne.jsx)(b.a,{children:o.slice(a*r,a*r+r).map((e=>Object(ne.jsx)(x.a,{hover:!0,role:"checkbox",tabIndex:-1,children:pe.map((a=>{const t=e[a.id];return Object(ne.jsx)(m.a,{align:a.align,children:a.format?a.format(t,e,l):t},a.id)}))},e.code)))})]})}),Object(ne.jsx)(p.a,{rowsPerPageOptions:[10,25,100],component:"div",count:o.length,rowsPerPage:r,page:a,onChangePage:(e,a)=>{t(a)},onChangeRowsPerPage:e=>{c(+e.target.value),t(0)}})]})}var ge=t(23),fe=t.n(ge),ve=t(214),we=t(222),ke=t(227),Ce=t(111),ye=t.n(Ce),Se=t(226),Re=t(224),Pe=t(223),Ie=t(225);const Ne=Object(d.a)((e=>({root:{width:"100%"}})));function Be(e){const a=Ne(),[t,c]=n.a.useState(!1),[o,i]=n.a.useState(e.dialog.title),[l,d]=n.a.useState(e.dialog.content);Object(r.useImperativeHandle)(e.onRef,(()=>({handleClickOpen:s})));const s=()=>{c(!0)},j=()=>{c(!1)};return Object(ne.jsx)("div",{className:a.root,children:Object(ne.jsxs)(S.a,{open:t,onClose:j,"aria-labelledby":"form-dialog-title",children:[Object(ne.jsx)(Pe.a,{id:"form-dialog-title",children:o}),Object(ne.jsx)(Re.a,{children:Object(ne.jsx)(Ie.a,{children:l})}),Object(ne.jsx)(Se.a,{children:Object(ne.jsx)(y.a,{onClick:j,color:"primary",children:"OK"})})]})})}const Te=Object(d.a)((e=>({appBar:{position:"relative"},title:{marginLeft:e.spacing(2),flex:1},formClass:{"& > *":{margin:e.spacing(1),width:"25ch"}},formDiv:{margin:"0 auto"}}))),Le=n.a.forwardRef((function(e,a){return Object(ne.jsx)(L.a,{direction:"up",ref:a,...e})}));function Ee(e){const a=Te(),[t,c]=n.a.useState(!1),[o,i]=n.a.useState(1===e.row.enabled),[l,d]=(e.row.cameraCode,n.a.useState({id:e.row.id,cameraId:e.row.cameraId,name:e.row.name,authCode:e.row.authCode,startTime:e.row.startTime?e.row.startTime:fe()(),deadline:e.row.deadline?e.row.deadline:fe()().add(7,"day"),enabled:e.row.enabled})),[s,j]=n.a.useState(!1),[b,m]=n.a.useState("");Object(r.useImperativeHandle)(e.onRef,(()=>({handleClickOpen:u})));const u=()=>{l.id=e.row.id,l.cameraId=e.row.cameraId,l.name=e.row.name,l.authCode=e.row.authCode,l.startTime=e.row.startTime?e.row.startTime:fe()(),l.deadline=e.row.deadline?e.row.deadline:fe()().add(7,"day"),l.enabled=e.row.enabled,d(l),c(!0)},h=()=>{c(!1)},p=a=>{ae(l).then((a=>{if(1===a.code)return c(!1),void(e.callBack&&e.callBack());m(a.msg),j(!0),window.setTimeout((function(){j(!1)}),5e3)}))},x=e=>{l[e.target.id]=fe()(e.target.value)};return Object(ne.jsx)("div",{children:Object(ne.jsxs)(S.a,{fullScreen:!0,open:t,onClose:h,TransitionComponent:Le,children:[Object(ne.jsx)(R.a,{className:a.appBar,children:Object(ne.jsxs)(P.a,{children:[Object(ne.jsx)(y.a,{variant:"contained",onClick:p,children:"\u4fdd\u5b58"}),Object(ne.jsx)("span",{children:"\xa0\xa0"}),"edit"===e.type?Object(ne.jsx)(y.a,{variant:"contained",onClick:a=>{te(l).then((a=>{if(1===a.code)return c(!1),void(e.callBack&&e.callBack());m(a.msg),j(!0),window.setTimeout((function(){j(!1)}),5e3)}))},children:"\u5220\u9664"}):"",Object(ne.jsx)(I.a,{variant:"h6",className:a.title,children:Object(ne.jsxs)("div",{children:[l.name?l.name+" \u5206\u4eab\u7f16\u8f91":"\u521b\u5efa\u5206\u4eab"," "]})}),Object(ne.jsx)(y.a,{autoFocus:!0,color:"inherit",onClick:h,children:Object(ne.jsx)(T.a,{})})]})}),s?Object(ne.jsxs)(de.a,{severity:"error",children:[Object(ne.jsx)(se.a,{children:"Error"}),b," ",Object(ne.jsx)("strong",{children:"check it out!"})]}):"",Object(ne.jsx)("div",{className:a.formDiv,children:Object(ne.jsxs)("form",{className:a.formClass,noValidate:!0,autoComplete:"off",onSubmit:p,children:["edit"===e.type?Object(ne.jsx)("div",{children:Object(ne.jsx)(le.a,{id:"id",label:"id",InputProps:{readOnly:!0},defaultValue:l.id})}):"",Object(ne.jsx)("div",{children:Object(ne.jsx)(le.a,{id:"name",label:"\u5206\u4eab\u540d\u79f0",defaultValue:l.name,onChange:e=>{l[e.target.id]=e.target.value}})}),Object(ne.jsx)("div",{children:Object(ne.jsx)(le.a,{id:"startTime",label:"\u5f00\u59cb\u65f6\u95f4",type:"datetime-local",defaultValue:l.startTime?fe()(l.startTime).format("YYYY-MM-DDTHH:mm"):fe()().format("YYYY-MM-DDTHH:mm"),InputLabelProps:{shrink:!0},onChange:x})}),Object(ne.jsx)("div",{children:Object(ne.jsx)(le.a,{id:"deadline",label:"\u622a\u6b62\u65e5\u671f",type:"datetime-local",defaultValue:l.deadline?fe()(l.deadline).format("YYYY-MM-DDTHH:mm"):fe()().add(7,"day").format("YYYY-MM-DDTHH:mm"),InputLabelProps:{shrink:!0},onChange:x})}),"edit"===e.type?"":Object(ne.jsx)("div",{children:Object(ne.jsx)(N.a,{checked:o,id:"enabled",onChange:e=>{console.log(e.target.checked),l[e.target.id]=e.target.checked?1:0,"enabled"===e.target.id&&i(e.target.checked)},color:"primary",name:"enabled",inputProps:{"aria-label":"primary checkbox"}})})]})})]})})}const Me=Object(d.a)((()=>({root:{position:"relative"},dropdown:{position:"absolute",top:28,right:0,left:0,zIndex:999,border:"0px solid",padding:1,backgroundColor:"#f8f8f8"}})));function De(e){const a=Me(),[t,r]=n.a.useState(!1),[c,o]=n.a.useState({title:"Success",content:"copy success !"}),[i,l]=n.a.useState(e.row);let d=n.a.createRef();let s=n.a.createRef();let j=n.a.createRef();return Object(ne.jsx)(w.a,{onClickAway:()=>{r(!1)},children:Object(ne.jsxs)("div",{className:a.root,children:[Object(ne.jsx)("button",{type:"button",onClick:()=>{r((e=>!e))},children:Object(ne.jsx)(C.a,{})}),t?Object(ne.jsx)("div",{className:a.dropdown,children:Object(ne.jsxs)(ve.a,{component:"nav","aria-label":"secondary mailbox folders",children:[Object(ne.jsx)(we.a,{button:!0,children:Object(ne.jsx)(ke.a,{primary:"\u7f16\u8f91",onClick:function(){r(!1),s.current.handleClickOpen()}})}),Object(ne.jsx)(we.a,{button:!0,onClick:()=>{i.enabled=1===i.enabled?0:1,re(i).then((a=>{if(1===a.code)return r(!1),void(e.callBack&&e.callBack())}))},children:1===e.row.enabled?Object(ne.jsx)(ke.a,{primary:"\u7981\u7528"}):Object(ne.jsx)(ke.a,{primary:"\u542f\u7528"})}),Object(ne.jsx)(we.a,{button:!0,onClick:function(){r(!1),l(e.row),d.current.handleClickOpen()},children:Object(ne.jsx)(ke.a,{primary:"\u64ad\u653e"})}),Object(ne.jsx)(we.a,{button:!0,onClick:function(){r(!1);let a=window.location.origin+window.location.pathname+"#/live?method=temp&code="+e.row.cameraCode+"&authCode="+e.row.authCode;ye()(a),j.current.handleClickOpen()},children:Object(ne.jsx)(ke.a,{primary:"\u5206\u4eab"})})]})}):null,Object(ne.jsx)(ie,{playParam:{playMethod:"temp",cameraCode:i.cameraCode,authCode:i.authCode},onRef:d}),Object(ne.jsx)(Ee,{row:e.row,type:"edit",callBack:e.callBack,onRef:s}),Object(ne.jsx)(Be,{dialog:c,onRef:j})]})})}const Ve=[{id:"id",label:"id",minWidth:170},{id:"name",label:"\u5206\u4eab\u540d\u79f0"},{id:"authCode",label:"\u6743\u9650\u7801"},{id:"startTime",label:"\u5f00\u59cb\u65f6\u95f4",format:e=>{(new Date).getTimezoneOffset();return e?fe()(e).format("YYYY-MM-DD HH:mm"):"---"}},{id:"deadline",label:"\u622a\u6b62\u65f6\u95f4",format:e=>{(new Date).getTimezoneOffset();return e?fe()(e).format("YYYY-MM-DD HH:mm"):"---"}},{id:"enabled",label:"\u542f\u7528\u72b6\u6001",format:e=>Object(ne.jsx)(N.a,{checked:1===e,id:"enabled",color:"primary",name:"enabled",inputProps:{"aria-label":"primary checkbox"}})},{id:"action",label:"\u64cd\u4f5c",format:(e,a,t,r)=>(a.cameraId=t.id,a.cameraCode=t.code,Object(ne.jsx)(De,{row:a,callBack:r}))}],Ye=Object(d.a)((e=>({root:{width:"100%",position:"relative"},container:{minHeight:400},dropdown:{position:"absolute",top:28,right:0,left:0,zIndex:0,border:"0px solid",padding:1,backgroundColor:"#f8f8f8"},appBar:{position:"relative"},title:{marginLeft:e.spacing(2),flex:1}})));function Ue(e){const a=Ye(),[t,r]=n.a.useState(0),[c,o]=n.a.useState(10),[l,d]=n.a.useState([]),[O,g]=n.a.useState({id:"",cameraId:"",cameraCode:"",name:"",authCode:"",enabled:1,startTime:"",deadline:""}),[f,v]=n.a.useState({}),w=()=>{let e=(e=>{let a=new RegExp("(^|&|\\?)"+e+"=([^&]*)(&|$)","i"),t=window.location.hash.substr(1).match(a);return null!=t?decodeURIComponent(t[2]):null})("cameraId");q({id:e}).then((e=>{if(1===e.code){let a=e.data;if(!a)return;v(a);let t=O;t.cameraId=a.id,t.cameraCode=a.code,g(t),ee({cameraId:a.id}).then((e=>{1===e.code&&(l.splice(0),l.push(...e.data.page),d([]),d(l))}))}}))};let k=n.a.createRef();return n.a.useEffect(w,[]),Object(ne.jsxs)(s.a,{className:a.root,children:[Object(ne.jsx)(R.a,{className:a.appBar,children:Object(ne.jsxs)(P.a,{children:[Object(ne.jsx)(y.a,{variant:"contained",onClick:function(){g({}),g(O),k.current.handleClickOpen()},children:"\u521b\u5efa"}),Object(ne.jsxs)(I.a,{variant:"h6",className:a.title,children:["\u6444\u50cf\u5934: ",f.code?f.code:"--"," \u5206\u4eab\u5217\u8868"]}),Object(ne.jsx)(y.a,{variant:"contained",children:Object(ne.jsx)(i.b,{exact:!0,to:{pathname:"/",state:{}},children:"\u8fd4\u56de"})})]})}),Object(ne.jsx)(Ee,{row:O,type:"add",callBack:w,onRef:k}),Object(ne.jsx)(u.a,{className:a.container,children:Object(ne.jsxs)(j.a,{stickyHeader:!0,"aria-label":"sticky table",children:[Object(ne.jsx)(h.a,{children:Object(ne.jsx)(x.a,{children:Ve.map((e=>Object(ne.jsx)(m.a,{align:e.align,style:{minWidth:e.minWidth},children:e.label},e.id)))})}),Object(ne.jsx)(b.a,{children:l.slice(t*c,t*c+c).map((e=>Object(ne.jsx)(x.a,{hover:!0,role:"checkbox",tabIndex:-1,children:Ve.map((a=>{const t=e[a.id];return Object(ne.jsx)(m.a,{align:a.align,children:a.format?a.format(t,e,f,w):t},a.id)}))},e.id)))})]})}),Object(ne.jsx)(p.a,{rowsPerPageOptions:[10,25,100],component:"div",count:l.length,rowsPerPage:c,page:t,onChangePage:(e,a)=>{r(a)},onChangeRowsPerPage:e=>{o(+e.target.value),r(0)}})]})}var He=t(240),Fe=t(230),Ae=t(231),We=t(238),ze=t(228),$e=t(232),qe=t(229),Je=t(112),_e=t.n(Je);function Ke(){return Object(ne.jsxs)(I.a,{variant:"body2",color:"textSecondary",align:"center",children:["Copyright \xa9 ",Object(ne.jsx)(ze.a,{color:"inherit",href:"https://material-ui.com/",children:"Your Website"})," ",(new Date).getFullYear(),"."]})}const Ge=Object(d.a)((e=>({root:{height:"100vh"},image:{backgroundImage:"url(https://source.unsplash.com/random)",backgroundRepeat:"no-repeat",backgroundColor:"light"===e.palette.type?e.palette.grey[50]:e.palette.grey[900],backgroundSize:"cover",backgroundPosition:"center"},paper:{margin:e.spacing(8,4),display:"flex",flexDirection:"column",alignItems:"center"},avatar:{margin:e.spacing(1),backgroundColor:e.palette.secondary.main},form:{width:"100%",marginTop:e.spacing(1)},submit:{margin:e.spacing(3,0,2)}})));function Qe(){const e=Ge(),[a,t]=n.a.useState({userName:"",password:""}),[r,c]=n.a.useState(!1),[o,i]=n.a.useState(""),l=e=>{a[e.target.id]=e.target.value};return n.a.useEffect((()=>{"true"===window.localStorage.getItem("tokenExpired")&&localStorage.setItem("tokenExpired","false")}),[]),Object(ne.jsxs)(qe.a,{container:!0,component:"main",className:e.root,children:[Object(ne.jsx)(Fe.a,{}),Object(ne.jsx)(qe.a,{item:!0,xs:!1,sm:4,md:7,className:e.image}),Object(ne.jsx)(qe.a,{item:!0,xs:12,sm:8,md:5,component:s.a,elevation:6,square:!0,children:Object(ne.jsxs)("div",{className:e.paper,children:[Object(ne.jsx)(He.a,{className:e.avatar,children:Object(ne.jsx)(_e.a,{})}),Object(ne.jsx)(I.a,{component:"h1",variant:"h5",children:"Sign in"}),r?Object(ne.jsxs)(de.a,{severity:"error",children:[Object(ne.jsx)(se.a,{children:"Error"}),o," ",Object(ne.jsx)("strong",{children:"check it out!"})]}):"",Object(ne.jsxs)("form",{className:e.form,noValidate:!0,children:[Object(ne.jsx)(le.a,{variant:"outlined",margin:"normal",required:!0,fullWidth:!0,id:"userName",label:"UserName",name:"userName",autoComplete:"userName",autoFocus:!0,onChange:l}),Object(ne.jsx)(le.a,{variant:"outlined",margin:"normal",required:!0,fullWidth:!0,name:"password",label:"Password",type:"password",id:"password",autoComplete:"current-password",onChange:l}),Object(ne.jsx)(Ae.a,{control:Object(ne.jsx)(We.a,{value:"remember",color:"primary"}),label:"Remember me"}),Object(ne.jsx)(y.a,{fullWidth:!0,variant:"contained",color:"primary",className:e.submit,onClick:e=>{z(a).then((e=>{if(1===e.code){let a=window.localStorage;return a.setItem("token",e.data.token),a.setItem("tokenExpired","false"),void(window.location.hash="#/")}i(e.msg),c(!0)}))},children:"Sign In"}),Object(ne.jsxs)(qe.a,{container:!0,children:[Object(ne.jsx)(qe.a,{item:!0,xs:!0,children:Object(ne.jsx)(ze.a,{href:"#",variant:"body2",children:"Forgot password?"})}),Object(ne.jsx)(qe.a,{item:!0,children:Object(ne.jsx)(ze.a,{href:"#",variant:"body2",children:"Don't have an account? Sign Up"})})]}),Object(ne.jsx)($e.a,{mt:5,children:Object(ne.jsx)(Ke,{})})]})]})})]})}class Xe extends n.a.Component{render(){return Object(ne.jsx)("h2",{children:"404"})}}const Ze=Object(d.a)((e=>({appBar:{position:"relative",backgroundColor:"#eebbaa"},title:{marginLeft:e.spacing(2),flex:1},videoContainer:{width:"90%",margin:"0 auto"}})));function ea(e){const a=Ze(),[t,r]=n.a.useState(!0);var c=null,o=0;const i=e=>{let a=new RegExp("(^|&|\\?)"+e+"=([^&]*)(&|$)","i"),t=window.location.hash.substr(1).match(a);return null!=t?decodeURIComponent(t[2]):null},l=e=>{let a=i("method"),r=i("code"),n=i("authCode");if(!a||!r||!n)return;var c={type:"flv"};let o=W+"/live/"+a+"/"+r+"/"+n+".flv";c.url=o,c.hasAudio=t,c.isLive=!0,console.log("MediaDataSource",c),d(c)},d=e=>{var a=document.getElementsByClassName("centeredVideo")[0];"undefined"!==typeof c&&null!=c&&(c.pause(),c.unload(),c.detachMediaElement(),c.destroy(),c=null),(c=M.a.createPlayer(e,{enableWorker:!1,lazyLoadMaxDuration:180,seekType:"range"})).on(M.a.Events.ERROR,((e,a,t)=>{console.log("errorType:",e),console.log("errorDetail:",a),console.log("errorInfo:",t),c&&(c.pause(),c.unload(),c.detachMediaElement(),c.destroy(),c=null,window.setTimeout(l,500))})),c.on(M.a.Events.STATISTICS_INFO,(function(e){0!=o?o!=e.decodedFrames?o=e.decodedFrames:(console.log("decodedFrames:",e.decodedFrames),o=0,c&&(c.pause(),c.unload(),c.detachMediaElement(),c.destroy(),c=null,window.setTimeout(l,500))):o=e.decodedFrames})),c.attachMediaElement(a),c.load(),c.play()};return Object(ne.jsx)("div",{children:Object(ne.jsxs)("div",{children:[Object(ne.jsx)(R.a,{className:a.appBar,children:Object(ne.jsxs)(P.a,{children:[Object(ne.jsx)(y.a,{variant:"contained",onClick:l,children:"play"}),Object(ne.jsxs)(I.a,{variant:"h6",className:a.title,children:["hasAudio",Object(ne.jsx)(N.a,{checked:t,id:"Audio",color:"primary",name:"hasAudio",onChange:e=>{r(e.target.checked)},inputProps:{"aria-label":"primary checkbox"}})]})]})}),Object(ne.jsx)("div",{className:a.videoContainer,children:Object(ne.jsx)("div",{children:Object(ne.jsx)("video",{name:"videoElement",className:"centeredVideo",controls:!0,allow:"autoPlay",width:"100%",children:"Your browser is too old which doesn't support HTML5 video."})})})]})})}function aa(e){return Object(ne.jsxs)(l.c,{children:[Object(ne.jsx)(l.a,{exact:!0,path:"/",component:Oe}),Object(ne.jsx)(l.a,{exact:!0,path:"/camerashare",component:Ue}),Object(ne.jsx)(l.a,{exact:!0,path:"/live",component:ea}),Object(ne.jsx)(l.a,{path:"/login",component:Qe}),Object(ne.jsx)(l.a,{component:Xe})]})}var ta=Object(r.memo)(aa);var ra=function(){return Object(ne.jsx)("div",{children:Object(ne.jsx)(i.a,{children:Object(ne.jsx)(ta,{})})})};function na(){return Object(ne.jsx)(ra,{})}o.a.render(Object(ne.jsx)(na,{}),document.querySelector("#app"))}},[[165,1,2]]]);
//# sourceMappingURL=main.12e4388c.chunk.js.map