// Code generated by go generate; DO NOT EDIT.

package static // import "miniflux.app/ui/static"

var Javascripts = map[string]string{
	"app": `(function(){'use strict';class DomHelper{static isVisible(element){return element.offsetParent!==null;}
static openNewTab(url){let win=window.open("");win.opener=null;win.location=url;win.focus();}
static scrollPageTo(element){let windowScrollPosition=window.pageYOffset;let windowHeight=document.documentElement.clientHeight;let viewportPosition=windowScrollPosition+windowHeight;let itemBottomPosition=element.offsetTop+element.offsetHeight;if(viewportPosition-itemBottomPosition<0||viewportPosition-element.offsetTop>windowHeight){window.scrollTo(0,element.offsetTop-10);}}
static getVisibleElements(selector){let elements=document.querySelectorAll(selector);let result=[];for(let i=0;i<elements.length;i++){if(this.isVisible(elements[i])){result.push(elements[i]);}}
return result;}
static findParent(element,selector){for(;element&&element!==document;element=element.parentNode){if(element.classList.contains(selector)){return element;}}
return null;}
static hasPassiveEventListenerOption(){var passiveSupported=false;try{var options=Object.defineProperty({},"passive",{get:function(){passiveSupported=true;}});window.addEventListener("test",options,options);window.removeEventListener("test",options,options);}catch(err){passiveSupported=false;}
return passiveSupported;}}
class TouchHandler{constructor(navHandler){this.navHandler=navHandler;this.reset();}
reset(){this.touch={start:{x:-1,y:-1},move:{x:-1,y:-1},element:null};}
calculateDistance(){if(this.touch.start.x>=-1&&this.touch.move.x>=-1){let horizontalDistance=Math.abs(this.touch.move.x-this.touch.start.x);let verticalDistance=Math.abs(this.touch.move.y-this.touch.start.y);if(horizontalDistance>30&&verticalDistance<70){return this.touch.move.x-this.touch.start.x;}}
return 0;}
findElement(element){if(element.classList.contains("touch-item")){return element;}
return DomHelper.findParent(element,"touch-item");}
onTouchStart(event){if(event.touches===undefined||event.touches.length!==1){return;}
this.reset();this.touch.start.x=event.touches[0].clientX;this.touch.start.y=event.touches[0].clientY;this.touch.element=this.findElement(event.touches[0].target);}
onTouchMove(event){if(event.touches===undefined||event.touches.length!==1||this.element===null){return;}
this.touch.move.x=event.touches[0].clientX;this.touch.move.y=event.touches[0].clientY;let distance=this.calculateDistance();let absDistance=Math.abs(distance);if(absDistance>0){let opacity=1-(absDistance>75?0.9:absDistance/75*0.9);let tx=distance>75?75:(distance<-75?-75:distance);this.touch.element.style.opacity=opacity;this.touch.element.style.transform="translateX("+tx+"px)";event.preventDefault();}}
onTouchEnd(event){if(event.touches===undefined){return;}
if(this.touch.element!==null){let distance=Math.abs(this.calculateDistance());if(distance>75){EntryHandler.toggleEntryStatus(this.touch.element);}
this.touch.element.style.opacity=1;this.touch.element.style.transform="none";}
this.reset();}
listen(){let elements=document.querySelectorAll(".touch-item");let hasPassiveOption=DomHelper.hasPassiveEventListenerOption();elements.forEach((element)=>{element.addEventListener("touchstart",(e)=>this.onTouchStart(e),hasPassiveOption?{passive:true}:false);element.addEventListener("touchmove",(e)=>this.onTouchMove(e),hasPassiveOption?{passive:false}:false);element.addEventListener("touchend",(e)=>this.onTouchEnd(e),hasPassiveOption?{passive:true}:false);element.addEventListener("touchcancel",()=>this.reset(),hasPassiveOption?{passive:true}:false);});let entryContentElement=document.querySelector(".entry-content");if(entryContentElement){let doubleTapTimers={previous:null,next:null};const detectDoubleTap=(doubleTapTimer,event)=>{const timer=doubleTapTimers[doubleTapTimer];if(timer===null){doubleTapTimers[doubleTapTimer]=setTimeout(()=>{doubleTapTimers[doubleTapTimer]=null;},200);}else{event.preventDefault();this.navHandler.goToPage(doubleTapTimer);}};entryContentElement.addEventListener("touchend",(e)=>{if(e.changedTouches[0].clientX>=(entryContentElement.offsetWidth/2)){detectDoubleTap("next",e);}else{detectDoubleTap("previous",e);}},hasPassiveOption?{passive:false}:false);entryContentElement.addEventListener("touchmove",(e)=>{Object.keys(doubleTapTimers).forEach(timer=>doubleTapTimers[timer]=null);});}}}
class KeyboardHandler{constructor(){this.queue=[];this.shortcuts={};}
on(combination,callback){this.shortcuts[combination]=callback;}
listen(){document.onkeydown=(event)=>{if(this.isEventIgnored(event)||this.isModifierKeyDown(event)){return;}
let key=this.getKey(event);this.queue.push(key);for(let combination in this.shortcuts){let keys=combination.split(" ");if(keys.every((value,index)=>value===this.queue[index])){this.queue=[];this.shortcuts[combination](event);return;}
if(keys.length===1&&key===keys[0]){this.queue=[];this.shortcuts[combination](event);return;}}
if(this.queue.length>=2){this.queue=[];}};}
isEventIgnored(event){return event.target.tagName==="INPUT"||event.target.tagName==="TEXTAREA";}
isModifierKeyDown(event){return event.getModifierState("Control")||event.getModifierState("Alt")||event.getModifierState("Meta");}
getKey(event){const mapping={'Esc':'Escape','Up':'ArrowUp','Down':'ArrowDown','Left':'ArrowLeft','Right':'ArrowRight'};for(let key in mapping){if(mapping.hasOwnProperty(key)&&key===event.key){return mapping[key];}}
return event.key;}}
class MouseHandler{onClick(selector,callback,noPreventDefault){let elements=document.querySelectorAll(selector);elements.forEach((element)=>{element.onclick=(event)=>{if(!noPreventDefault){event.preventDefault();}
callback(event);};});}}class FormHandler{static handleSubmitButtons(){let elements=document.querySelectorAll("form");elements.forEach((element)=>{element.onsubmit=()=>{let button=document.querySelector("button");if(button){button.innerHTML=button.dataset.labelLoading;button.disabled=true;}};});}}
class RequestBuilder{constructor(url){this.callback=null;this.url=url;this.options={method:"POST",cache:"no-cache",credentials:"include",body:null,headers:new Headers({"Content-Type":"application/json","X-Csrf-Token":this.getCsrfToken()})};}
withBody(body){this.options.body=JSON.stringify(body);return this;}
withCallback(callback){this.callback=callback;return this;}
getCsrfToken(){let element=document.querySelector("meta[name=X-CSRF-Token]");if(element!==null){return element.getAttribute("value");}
return "";}
execute(){fetch(new Request(this.url,this.options)).then((response)=>{if(this.callback){this.callback(response);}});}}
class UnreadCounterHandler{static decrement(n){this.updateValue((current)=>{return current-n;});}
static increment(n){this.updateValue((current)=>{return current+n;});}
static updateValue(callback){let counterElements=document.querySelectorAll("span.unread-counter");counterElements.forEach((element)=>{let oldValue=parseInt(element.textContent,10);element.innerHTML=callback(oldValue);});if(window.location.href.endsWith('/unread')){let oldValue=parseInt(document.title.split('(')[1],10);let newValue=callback(oldValue);document.title=document.title.replace(/(.*?)\(\d+\)(.*?)/,function(match,prefix,suffix,offset,string){return prefix+'('+newValue+')'+suffix;});}}}
class EntryHandler{static updateEntriesStatus(entryIDs,status,callback){let url=document.body.dataset.entriesStatusUrl;let request=new RequestBuilder(url);request.withBody({entry_ids:entryIDs,status:status});request.withCallback(callback);request.execute();if(status==="read"){UnreadCounterHandler.decrement(1);}else{UnreadCounterHandler.increment(1);}}
static toggleEntryStatus(element,silently=false){let entryID=parseInt(element.dataset.id,10);let link=element.querySelector("a[data-toggle-status]");let currentStatus=link.dataset.value;let newStatus=currentStatus==="read"?"unread":"read";this.updateEntriesStatus([entryID],newStatus);if(currentStatus==="read"){link.innerHTML=link.dataset.labelRead;link.dataset.value="unread";}else{link.innerHTML=link.dataset.labelUnread;link.dataset.value="read";}
if(element.classList.contains("item-status-"+currentStatus)&&!silently){element.classList.remove("item-status-"+currentStatus);element.classList.add("item-status-"+newStatus);}}
static toggleBookmark(element){element.innerHTML=element.dataset.labelLoading;let request=new RequestBuilder(element.dataset.bookmarkUrl);request.withCallback(()=>{if(element.dataset.value==="star"){element.innerHTML=element.dataset.labelStar;element.dataset.value="unstar";}else{element.innerHTML=element.dataset.labelUnstar;element.dataset.value="star";}});request.execute();}
static markEntryAsRead(element){if(element.classList.contains("item-status-unread")){element.classList.remove("item-status-unread");element.classList.add("item-status-read");let entryID=parseInt(element.dataset.id,10);this.updateEntriesStatus([entryID],"read");}}
static saveEntry(element){if(element.dataset.completed){return;}
element.innerHTML=element.dataset.labelLoading;let request=new RequestBuilder(element.dataset.saveUrl);request.withCallback(()=>{element.innerHTML=element.dataset.labelDone;element.dataset.completed=true;});request.execute();}
static fetchOriginalContent(element){if(element.dataset.completed){return;}
element.innerHTML=element.dataset.labelLoading;let request=new RequestBuilder(element.dataset.fetchContentUrl);request.withCallback((response)=>{element.innerHTML=element.dataset.labelDone;element.dataset.completed=true;response.json().then((data)=>{if(data.hasOwnProperty("content")){document.querySelector(".entry-content").innerHTML=data.content;}});});request.execute();}}
class FeedHandler{static unsubscribe(feedUrl,callback){let request=new RequestBuilder(feedUrl);request.withCallback(callback);request.execute();}}
class ConfirmHandler{executeRequest(url,redirectURL){let request=new RequestBuilder(url);request.withCallback(()=>{if(redirectURL){window.location.href=redirectURL;}else{window.location.reload();}});request.execute();}
handle(event){let questionElement=document.createElement("span");let linkElement=event.target;let containerElement=linkElement.parentNode;linkElement.style.display="none";let yesElement=document.createElement("a");yesElement.href="#";yesElement.appendChild(document.createTextNode(linkElement.dataset.labelYes));yesElement.onclick=(event)=>{event.preventDefault();let loadingElement=document.createElement("span");loadingElement.className="loading";loadingElement.appendChild(document.createTextNode(linkElement.dataset.labelLoading));questionElement.remove();containerElement.appendChild(loadingElement);this.executeRequest(linkElement.dataset.url,linkElement.dataset.redirectUrl);};let noElement=document.createElement("a");noElement.href="#";noElement.appendChild(document.createTextNode(linkElement.dataset.labelNo));noElement.onclick=(event)=>{event.preventDefault();linkElement.style.display="inline";questionElement.remove();};questionElement.className="confirm";questionElement.appendChild(document.createTextNode(linkElement.dataset.labelQuestion+" "));questionElement.appendChild(yesElement);questionElement.appendChild(document.createTextNode(", "));questionElement.appendChild(noElement);containerElement.appendChild(questionElement);}}
class MenuHandler{clickMenuListItem(event){let element=event.target;if(element.tagName==="A"){window.location.href=element.getAttribute("href");}else{window.location.href=element.querySelector("a").getAttribute("href");}}
toggleMainMenu(){let menu=document.querySelector(".header nav ul");if(DomHelper.isVisible(menu)){menu.style.display="none";}else{menu.style.display="block";}
let searchElement=document.querySelector(".header .search");if(DomHelper.isVisible(searchElement)){searchElement.style.display="none";}else{searchElement.style.display="block";}}}
class ModalHandler{static exists(){return document.getElementById("modal-container")!==null;}
static open(fragment){if(ModalHandler.exists()){return;}
let container=document.createElement("div");container.id="modal-container";container.appendChild(document.importNode(fragment,true));document.body.appendChild(container);let closeButton=document.querySelector("a.btn-close-modal");if(closeButton!==null){closeButton.onclick=(event)=>{event.preventDefault();ModalHandler.close();};}}
static close(){let container=document.getElementById("modal-container");if(container!==null){container.parentNode.removeChild(container);}}}
class NavHandler{setFocusToSearchInput(event){event.preventDefault();event.stopPropagation();let toggleSwitchElement=document.querySelector(".search-toggle-switch");if(toggleSwitchElement){toggleSwitchElement.style.display="none";}
let searchFormElement=document.querySelector(".search-form");if(searchFormElement){searchFormElement.style.display="block";}
let searchInputElement=document.getElementById("search-input");if(searchInputElement){searchInputElement.focus();searchInputElement.value="";}}
showKeyboardShortcuts(){let template=document.getElementById("keyboard-shortcuts");if(template!==null){ModalHandler.open(template.content);}}
markPageAsRead(showOnlyUnread){let items=DomHelper.getVisibleElements(".items .item");let entryIDs=[];items.forEach((element)=>{element.classList.add("item-status-read");entryIDs.push(parseInt(element.dataset.id,10));});if(entryIDs.length>0){EntryHandler.updateEntriesStatus(entryIDs,"read",()=>{if(showOnlyUnread){window.location.reload();}else{this.goToPage("next",true);}});}}
saveEntry(){if(this.isListView()){let currentItem=document.querySelector(".current-item");if(currentItem!==null){let saveLink=currentItem.querySelector("a[data-save-entry]");if(saveLink){EntryHandler.saveEntry(saveLink);}}}else{let saveLink=document.querySelector("a[data-save-entry]");if(saveLink){EntryHandler.saveEntry(saveLink);}}}
fetchOriginalContent(){if(!this.isListView()){let link=document.querySelector("a[data-fetch-content-entry]");if(link){EntryHandler.fetchOriginalContent(link);}}}
toggleEntryStatus(){if(!this.isListView()){EntryHandler.toggleEntryStatus(document.querySelector(".entry"));return;}
let currentItem=document.querySelector(".current-item");if(currentItem!==null){this.goToNextListItem();EntryHandler.toggleEntryStatus(currentItem);}}
toggleBookmark(){if(!this.isListView()){this.toggleBookmarkLink(document.querySelector(".entry"));return;}
let currentItem=document.querySelector(".current-item");if(currentItem!==null){this.toggleBookmarkLink(currentItem);}}
toggleBookmarkLink(parent){let bookmarkLink=parent.querySelector("a[data-toggle-bookmark]");if(bookmarkLink){EntryHandler.toggleBookmark(bookmarkLink);}}
openOriginalLink(){let entryLink=document.querySelector(".entry h1 a");if(entryLink!==null){DomHelper.openNewTab(entryLink.getAttribute("href"));return;}
let currentItemOriginalLink=document.querySelector(".current-item a[data-original-link]");if(currentItemOriginalLink!==null){DomHelper.openNewTab(currentItemOriginalLink.getAttribute("href"));let currentItem=document.querySelector(".current-item");this.goToNextListItem();EntryHandler.markEntryAsRead(currentItem);}}
openSelectedItem(){let currentItemLink=document.querySelector(".current-item .item-title a");if(currentItemLink!==null){window.location.href=currentItemLink.getAttribute("href");}}
unsubscribeFromFeed(){let unsubscribeLinks=document.querySelectorAll("[data-action=remove-feed]");if(unsubscribeLinks.length===1){let unsubscribeLink=unsubscribeLinks[0];FeedHandler.unsubscribe(unsubscribeLink.dataset.url,()=>{if(unsubscribeLink.dataset.redirectUrl){window.location.href=unsubscribeLink.dataset.redirectUrl;}else{window.location.reload();}});}}
goToPage(page,fallbackSelf){let element=document.querySelector("a[data-page="+page+"]");if(element){document.location.href=element.href;}else if(fallbackSelf){window.location.reload();}}
goToPrevious(){if(this.isListView()){this.goToPreviousListItem();}else{this.goToPage("previous");}}
goToNext(){if(this.isListView()){this.goToNextListItem();}else{this.goToPage("next");}}
goToFeedOrFeeds(){if(this.isEntry()){let feedAnchor=document.querySelector("span.entry-website a");if(feedAnchor!==null){window.location.href=feedAnchor.href;}}else{this.goToPage('feeds');}}
goToPreviousListItem(){let items=DomHelper.getVisibleElements(".items .item");if(items.length===0){return;}
if(document.querySelector(".current-item")===null){items[0].classList.add("current-item");items[0].querySelector('.item-header a').focus();return;}
for(let i=0;i<items.length;i++){if(items[i].classList.contains("current-item")){items[i].classList.remove("current-item");if(i-1>=0){items[i-1].classList.add("current-item");DomHelper.scrollPageTo(items[i-1]);items[i-1].querySelector('.item-header a').focus();}
break;}}}
goToNextListItem(){let currentItem=document.querySelector(".current-item");let items=DomHelper.getVisibleElements(".items .item");if(items.length===0){return;}
if(currentItem===null){items[0].classList.add("current-item");items[0].querySelector('.item-header a').focus();return;}
for(let i=0;i<items.length;i++){if(items[i].classList.contains("current-item")){items[i].classList.remove("current-item");if(i+1<items.length){items[i+1].classList.add("current-item");DomHelper.scrollPageTo(items[i+1]);items[i+1].querySelector('.item-header a').focus();}
break;}}}
isEntry(){return document.querySelector("section.entry")!==null;}
isListView(){return document.querySelector(".items")!==null;}}
class LinkStateHandler{static flip(element){let labelElement=document.createElement("span");labelElement.className="link-flipped-state";labelElement.appendChild(document.createTextNode(element.dataset.labelNewState));element.parentNode.appendChild(labelElement);element.parentNode.removeChild(element);}}
class ArticleHandler{static load(element){let elements=element.querySelectorAll(".article_view_url");elements.forEach((element)=>{let loadingElementWrapper=document.createElement("div");loadingElementWrapper.className="lds-dual-ring-wrapper";let loadingElement=document.createElement("div");loadingElement.className="lds-dual-ring";loadingElementWrapper.appendChild(loadingElement)
element.parentNode.appendChild(loadingElementWrapper);let request=new RequestBuilder(element.href);request.withCallback((data)=>{data.json().then(function(json){let view=document.createElement("div");view.className="entry-content";view.innerHTML=json.content;loadingElementWrapper.remove();element.parentNode.appendChild(view);element.remove();});});request.execute();});}}
class AppearHandler{constructor(selector="",opts={}){this.selectors=[];this.checkBinded=false;this.checkLock=false;this.defaults={interval:250,force_process:false};this.priorAppeared=[];if(selector!=""){this.addSelector(selector,opts)}
this.startMonitorLoop();}
process(){this.checkLock=false;function isVisible(element){return!!(element.offsetWidth||element.offsetHeight||element.getClientRects().length);}
function isAppeared(element){if(!isVisible(element)){return false;}
let windowLeft=window.scrollX;let windowTop=window.scrollY;let left=element.offsetLeft;let top=element.offsetTop;let belowTopEdge=top+element.clientHeight>=windowTop;let aboveBottomEdge=top<=windowTop+window.innerHeight;let rightAfterLeftEdge=left+element.clientWidth>=windowLeft;let leftBeforeRightEdge=left<=windowLeft+window.innerWidth;element.dataset.belowTopEdge=belowTopEdge;return belowTopEdge&&aboveBottomEdge&&rightAfterLeftEdge&&leftBeforeRightEdge;}
for(var index=0,selectorsLength=this.selectors.length;index<selectorsLength;index++){var appeared=this.selectors[index].filter(isAppeared);appeared.filter(element=>element.dataset._appearTriggered!="true").map(element=>element.dispatchEvent(new Event("appear")));if(this.priorAppeared[index]){var disappeared=this.priorAppeared[index].filter(x=>!appeared.includes(x));disappeared.filter(element=>element.dataset._appearTriggered=="true").map(element=>element.dispatchEvent(new Event("disappear")));}
this.priorAppeared[index]=appeared;}}
addSelector(selector,opts={}){if(typeof(opts.onappear)=="undefined"){opts.onappear=function(e){};}
if(typeof(opts.ondisappear)=="undefined"){opts.ondisappear=function(e){};}
let elements=Array.prototype.slice.call(document.querySelectorAll(selector));this.priorAppeared.push(elements);this.selectors.push(elements);elements.map(element=>element.addEventListener('appear',function(e){element=e.target;element.dataset._appearTriggered="true";opts.onappear(element);}));elements.map(element=>element.addEventListener('disappear',function(e){element=e.target;element.dataset._appearTriggered="false";opts.ondisappear(element);}));this.checkAfterFewSeconds();}
checkAfterFewSeconds(){if(this.checkLock){return;}
this.checkLock=true;setTimeout(this.process.bind(this),this.defaults.interval);}
startMonitorLoop(){if(!this.checkBinded){window.addEventListener('scroll',this.checkAfterFewSeconds.bind(this));window.addEventListener('resize',this.checkAfterFewSeconds.bind(this));this.checkBinded=true;}}}
document.addEventListener("DOMContentLoaded",function(){FormHandler.handleSubmitButtons();let navHandler=new NavHandler();if(!document.querySelector("body[data-disable-keyboard-shortcuts=true]")){let keyboardHandler=new KeyboardHandler();keyboardHandler.on("g u",()=>navHandler.goToPage("unread"));keyboardHandler.on("g b",()=>navHandler.goToPage("starred"));keyboardHandler.on("g h",()=>navHandler.goToPage("history"));keyboardHandler.on("g f",()=>navHandler.goToFeedOrFeeds());keyboardHandler.on("g c",()=>navHandler.goToPage("categories"));keyboardHandler.on("g s",()=>navHandler.goToPage("settings"));keyboardHandler.on("ArrowLeft",()=>navHandler.goToPrevious());keyboardHandler.on("ArrowRight",()=>navHandler.goToNext());keyboardHandler.on("k",()=>navHandler.goToPrevious());keyboardHandler.on("p",()=>navHandler.goToPrevious());keyboardHandler.on("j",()=>navHandler.goToNext());keyboardHandler.on("n",()=>navHandler.goToNext());keyboardHandler.on("h",()=>navHandler.goToPage("previous"));keyboardHandler.on("l",()=>navHandler.goToPage("next"));keyboardHandler.on("o",()=>navHandler.openSelectedItem());keyboardHandler.on("v",()=>navHandler.openOriginalLink());keyboardHandler.on("m",()=>navHandler.toggleEntryStatus());keyboardHandler.on("A",()=>{let element=document.querySelector("a[data-on-click=markPageAsRead]");navHandler.markPageAsRead(element.dataset.showOnlyUnread||false);});keyboardHandler.on("s",()=>navHandler.saveEntry());keyboardHandler.on("d",()=>navHandler.fetchOriginalContent());keyboardHandler.on("f",()=>navHandler.toggleBookmark());keyboardHandler.on("?",()=>navHandler.showKeyboardShortcuts());keyboardHandler.on("#",()=>navHandler.unsubscribeFromFeed());keyboardHandler.on("/",(e)=>navHandler.setFocusToSearchInput(e));keyboardHandler.on("Escape",()=>ModalHandler.close());keyboardHandler.listen();}
let touchHandler=new TouchHandler(navHandler);touchHandler.listen();let appearHandler=new AppearHandler();appearHandler.addSelector(".item-status-unread",{"onappear":function(element){if(document.querySelector("body[data-entry-embedded=true]")){ArticleHandler.load(element);}}});appearHandler.addSelector(".item-header",{"ondisappear":function(element){if(!document.querySelector("body[data-auto-mark-as-read=true]")){return;}
let currentStatus=element.parentNode.querySelector("a[data-toggle-status]").dataset.value;if(element.dataset.belowTopEdge=="false"&&currentStatus=="unread"){EntryHandler.toggleEntryStatus(element.parentNode,"silently");}}});let mouseHandler=new MouseHandler();mouseHandler.onClick("a[data-save-entry]",(event)=>{EntryHandler.saveEntry(event.target);});mouseHandler.onClick("a[data-toggle-bookmark]",(event)=>{EntryHandler.toggleBookmark(event.target);});mouseHandler.onClick("a[data-toggle-status]",(event)=>{let currentItem=DomHelper.findParent(event.target,"entry");if(!currentItem){currentItem=DomHelper.findParent(event.target,"item");}
if(currentItem){EntryHandler.toggleEntryStatus(currentItem);}});mouseHandler.onClick("a[data-fetch-content-entry]",(event)=>{EntryHandler.fetchOriginalContent(event.target);});mouseHandler.onClick("a[data-on-click=markPageAsRead]",(event)=>{navHandler.markPageAsRead(event.target.dataset.showOnlyUnread||false);});mouseHandler.onClick("a[data-confirm]",(event)=>{(new ConfirmHandler()).handle(event);});mouseHandler.onClick("a[data-action=search]",(event)=>{navHandler.setFocusToSearchInput(event);});mouseHandler.onClick("a[data-link-state=flip]",(event)=>{LinkStateHandler.flip(event.target);},true);if(document.documentElement.clientWidth<600){let menuHandler=new MenuHandler();mouseHandler.onClick(".logo",()=>menuHandler.toggleMainMenu());mouseHandler.onClick(".header nav li",(event)=>menuHandler.clickMenuListItem(event));}
mouseHandler.onClick(".toggle-entry-embedded",(event)=>{if(document.querySelector("body").dataset.entryEmbedded=="true"){document.querySelector("body").dataset.entryEmbedded="false";event.target.innerText="off";}else if(document.querySelector("body").dataset.entryEmbedded=="false"){document.querySelector("body").dataset.entryEmbedded="true";event.target.innerText="on";}
return false;});if("serviceWorker"in navigator){let scriptElement=document.getElementById("service-worker-script");if(scriptElement){navigator.serviceWorker.register(scriptElement.src);}}});})();`,
	"sw": `'use strict';self.addEventListener("fetch",(event)=>{if(event.request.url.includes("/feed/icon/")){event.respondWith(caches.open("feed_icons").then((cache)=>{return cache.match(event.request).then((response)=>{return response||fetch(event.request).then((response)=>{cache.put(event.request,response.clone());return response;});});}));}});`,
}

var JavascriptsChecksums = map[string]string{
	"app": "d489ac7a7750b02535edf40cfcb643416785f12960bf8df581b748732296a2c3",
	"sw":  "55fffa223919cc18572788fb9c62fccf92166c0eb5d3a1d6f91c31f24d020be9",
}
