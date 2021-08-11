(this.webpackJsonp=this.webpackJsonp||[]).push([["swag-test-environment"],{tIKo:function(e,t,n){"use strict";function r(e){return(r="function"==typeof Symbol&&"symbol"==typeof Symbol.iterator?function(e){return typeof e}:function(e){return e&&"function"==typeof Symbol&&e.constructor===Symbol&&e!==Symbol.prototype?"symbol":typeof e})(e)}function o(e,t){if(!(e instanceof t))throw new TypeError("Cannot call a class as a function")}function i(e,t){for(var n=0;n<t.length;n++){var r=t[n];r.enumerable=r.enumerable||!1,r.configurable=!0,"value"in r&&(r.writable=!0),Object.defineProperty(e,r.key,r)}}function c(e,t){return(c=Object.setPrototypeOf||function(e,t){return e.__proto__=t,e})(e,t)}function a(e){var t=function(){if("undefined"==typeof Reflect||!Reflect.construct)return!1;if(Reflect.construct.sham)return!1;if("function"==typeof Proxy)return!0;try{return Boolean.prototype.valueOf.call(Reflect.construct(Boolean,[],(function(){}))),!0}catch(e){return!1}}();return function(){var n,r=u(e);if(t){var o=u(this).constructor;n=Reflect.construct(r,arguments,o)}else n=r.apply(this,arguments);return s(this,n)}}function s(e,t){return!t||"object"!==r(t)&&"function"!=typeof t?function(e){if(void 0===e)throw new ReferenceError("this hasn't been initialised - super() hasn't been called");return e}(e):t}function u(e){return(u=Object.setPrototypeOf?Object.getPrototypeOf:function(e){return e.__proto__||Object.getPrototypeOf(e)})(e)}n.r(t);var l=Shopware.Classes.ApiService,f=function(e){!function(e,t){if("function"!=typeof t&&null!==t)throw new TypeError("Super expression must either be null or a function");e.prototype=Object.create(t&&t.prototype,{constructor:{value:e,writable:!0,configurable:!0}}),t&&c(e,t)}(u,e);var t,n,r,s=a(u);function u(e,t){var n=arguments.length>2&&void 0!==arguments[2]?arguments[2]:"_action/code-editor";return o(this,u),s.call(this,e,t,n)}return t=u,(n=[{key:"getFiles",value:function(){var e="".concat(this.getApiBasePath(),"/files");return this.httpClient.get(e,{headers:this.getBasicHeaders()}).then((function(e){return l.handleResponse(e)}))}},{key:"getFile",value:function(e){var t="".concat(this.getApiBasePath(),"/file");return this.httpClient.get(t,{params:{file:e},headers:this.getBasicHeaders()}).then((function(e){return e.data.content}))}},{key:"saveFile",value:function(e,t){var n="".concat(this.getApiBasePath(),"/file");return this.httpClient.put(n,{content:t},{params:{file:e},headers:this.getBasicHeaders()}).then((function(e){return l.handleResponse(e)}))}}])&&i(t.prototype,n),r&&i(t,r),u}(l),p=n("vcSJ"),d=n.n(p);function h(e,t,n,r,o,i,c){try{var a=e[i](c),s=a.value}catch(e){return void n(e)}a.done?t(s):Promise.resolve(s).then(r,o)}function v(e){return function(){var t=this,n=arguments;return new Promise((function(r,o){var i=e.apply(t,n);function c(e){h(i,r,o,c,a,"next",e)}function a(e){h(i,r,o,c,a,"throw",e)}c(void 0)}))}}var g=Shopware,m=g.Component,y=g.Mixin;m.register("code-editor-index",{inject:["codeEditorApiService"],template:d.a,mixins:[y.getByName("notification")],shortcuts:{"SYSTEMKEY+S":"saveFile"},data:function(){return{selectedLogFile:null,files:[],fileContent:""}},created:function(){var e=this;return v(regeneratorRuntime.mark((function t(){return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return t.next=2,e.codeEditorApiService.getFiles();case 2:e.files=t.sent;case 3:case"end":return t.stop()}}),t)})))()},methods:{onFileSelected:function(){var e=this;return v(regeneratorRuntime.mark((function t(){return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return e.fileContent="",t.next=3,e.codeEditorApiService.getFile(e.selectedLogFile);case 3:e.fileContent=t.sent;case 4:case"end":return t.stop()}}),t)})))()},saveFile:function(){var e=this;return v(regeneratorRuntime.mark((function t(){return regeneratorRuntime.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:if(null!==e.selectedLogFile){t.next=2;break}return t.abrupt("return");case 2:return t.next=4,e.codeEditorApiService.saveFile(e.selectedLogFile,e.fileContent);case 4:e.createNotificationSuccess({message:"File saved successfully"});case 5:case"end":return t.stop()}}),t)})))()}}}),Shopware.Module.register("code-editor",{type:"plugin",name:"code-editor.title",title:"code-editor.title",description:"",color:"#303A4F",icon:"default-device-dashboard",routes:{index:{component:"code-editor-index",path:"index"}},settingsItem:[{group:"plugins",to:"code.editor.index",icon:"default-action-settings",name:"code-editor.title"}]});var w=Shopware.Application;w.addServiceProvider("codeEditorApiService",(function(e){var t=w.getContainer("init");return new f(t.httpClient,e.loginService)}))},vcSJ:function(e,t){e.exports='<sw-page class="code-editor">\n    <template slot="content">\n        <sw-card-view>\n            <sw-card title="Code Editor">\n                <sw-alert variant="info">\n                    You can also use ALT+S to save the file\n                </sw-alert>\n\n                <label>Choose file:</label>\n                <sw-single-select\n                    :options="files"\n                    labelProperty="name"\n                    valueProperty="name"\n                    v-model="selectedLogFile"\n                    @change="onFileSelected"\n                ></sw-single-select>\n\n                <sw-code-editor :sanitizeInput="false" v-model="fileContent" v-if="fileContent"></sw-code-editor>\n\n                <sw-button variant="primary" @click="saveFile">Save file</sw-button>\n            </sw-card>\n        </sw-card-view>\n    </template>\n</sw-page>\n'}},[["tIKo","runtime"]]]);