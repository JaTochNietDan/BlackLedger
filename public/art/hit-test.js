// Hit testing uses the same building depth and sprite silhouette as rendering.
export function pickBuilding(buildings,masks,x,y){
 return [...buildings].sort((a,b)=>b.depth-a.depth).find(b=>{
  const u=(x-b.x)/b.w,v=(y-b.y)/b.h;
  if(u<0||u>=1||v<0||v>=1)return false;
  const mask=masks[b.id];
  return mask&&mask.data[(Math.floor(v*mask.height)*mask.width+Math.floor(u*mask.width))*4+3]>32;
 });
}
