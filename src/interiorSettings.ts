import type {InteriorPlace} from './interiorStaging';
type RoomSettings={lampIntensity?:number;name:string;model:string;span:number;leftWall:number;backWall:number;lamps:readonly (readonly [number,number,number])[]};
// Coordinates are glTF metres. Cutaway thresholds follow each authored wall;
// fixtures are positioned in the room that owns them, never borrowed by fallback.
export const interiorSettings:Record<InteriorPlace,RoomSettings>={
 chapel:{lampIntensity:1.8,name:'Thorne & Sons',model:'interior-chapel',span:7,leftWall:-5,backWall:-5.5,lamps:[[-2.1,2.9,.4],[2.1,2.9,.4]]},
 herald:{lampIntensity:2,name:'The Bellwether Herald',model:'interior-herald',span:7,leftWall:-5,backWall:-5,lamps:[[-2.5,2.98,0],[2.5,2.98,0]]},
 poolhall:{lampIntensity:2,name:'The Green Baize',model:'interior-poolhall',span:11,leftWall:-8,backWall:-9,lamps:[[-3.7,2.9,-4.25],[3.7,2.9,-4.25],[-3.7,2.9,0],[3.7,2.9,0],[-3.7,2.9,4.25],[3.7,2.9,4.25]]},
 tailor:{lampIntensity:2,name:'Ruttledge & Vance',model:'interior-tailor',span:7,leftWall:-5,backWall:-5,lamps:[[-2,2.82,-.6],[1.3,2.82,-.6]]},
 pawn:{lampIntensity:2,name:'Ackerman & Son',model:'interior-pawn',span:6,leftWall:-4,backWall:-4.5,lamps:[[-2,2.77,0],[1.9,2.77,0]]},
 riverside:{lampIntensity:2,name:'Riverside Courts',model:'interior-riverside',span:7,leftWall:-5,backWall:-5.5,lamps:[[-2,2.94,-1.9],[2.6,2.94,-1.4],[0,2.94,2.1]]},
 precinct:{lampIntensity:2,name:'Ward Street Station',model:'interior-precinct',span:7,leftWall:-5,backWall:-5.5,lamps:[[-2,2.94,-1.9],[2.6,2.94,-1.4],[0,2.94,2.1]]},
 market:{lampIntensity:2,name:'Mercer Exchange',model:'interior-exchange',span:7,leftWall:-5,backWall:-5,lamps:[[-2.5,2.97,.5],[1.4,2.97,.5],[1.4,2.97,-2.6]]},
 restaurant:{lampIntensity:3,name:"Vittoria’s",model:"interior-restaurant",span:6.5,leftWall:-5,backWall:-5.5,lamps:[[-2.7,2.68,-1.5],[-2.7,2.68,1.5],[2.7,2.68,-1.5],[2.7,2.68,1.5]]},
 lodging:{lampIntensity:1.4,name:'Your room at the Mariner',model:'interior-lodging-room',span:4.9,leftWall:-3,backWall:-3,lamps:[[-.67,1.3,-2.47]]},
 garage:{name:'Russo Motor Works',model:'interior-garage',span:8,leftWall:-6,backWall:-5.5,lamps:[[-3,3.47,-.5],[2.5,3.47,-.5]]},
 bar:{name:'Saint Agnes',model:'interior-saint-agnes',span:7,leftWall:-5.8,backWall:4.8,lamps:[[-3,2.65,2.7],[1,2.65,2.7],[4,2.65,2.7]]},
 mercercourt:{name:'Mercer Court',model:'interior-mercer-court',span:8.5,leftWall:-5.8,backWall:-6.8,lamps:[[-4.7,3.78,-6.2],[.1,3.78,-6.2]]},
 room:{name:'The Mariner',model:'interior-mariner',span:7.5,leftWall:-4.8,backWall:-5.8,lamps:[[-4.1,3.27,-5.46],[.5,3.27,-5.46]]},
 laundry:{name:'Bluebird Laundry',model:'interior-laundry',span:7.5,leftWall:-4.8,backWall:-5.8,lamps:[[-2.5,3.08,-1.5],[2.5,3.08,-1.5]]},
 estate:{name:'Cypress House',model:'interior-cypress',span:6.5,leftWall:-5,backWall:-4.5,lamps:[[-3.7,2,-2.7],[1,2,-2.7]]},
 apartment:{name:'Ashbury Court',model:'interior-ashbury',span:7,leftWall:-5,backWall:-5,lamps:[[-4.43,3.1,2.3],[-4.43,3.1,-1.1]]},
 flat:{name:'Your apartment',model:'interior-flat',span:5.7,leftWall:-4,backWall:-3.5,lamps:[[-2,2.8,-1],[1,2.8,-1]]},
 butcher:{name:'Fassano Meats',model:'interior-butcher',span:6.5,leftWall:-4.5,backWall:-4.5,lamps:[[-2.4,2.95,-1.5],[1.2,2.95,-1.5]]},
};
export function hasInterior(place:string):place is InteriorPlace {
 return Object.prototype.hasOwnProperty.call(interiorSettings,place);
}
