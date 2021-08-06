import { Injectable, OnInit, OnDestroy } from '@angular/core';

import { Observable } from 'rxjs';
import { RoomBackendService, CreateResponse, JoinResponse } from '../roombackend-service/roombackend.service';

export interface RoomData {
  roomName: string;
  previousMessages: string[];
}

@Injectable({
  providedIn: 'root'
})
export class RoomService implements OnInit, OnDestroy {

  constructor(private roomBackendService: RoomBackendService) {

  }

  ngOnDestroy() {

  }

  roomAuthenticationData(): RoomData | undefined {
    return this.roomAuthenticationInfo;
  }

  ngOnInit() {
  }

  create(roomName:string): Observable<CreateResponse> | undefined {
    const obs = this.roomBackendService.create(roomName);
    if (obs) {
      obs.subscribe({
        next: (resp:CreateResponse) => {
          this.roomAuthenticationInfo = { roomName: resp.roomName, previousMessages: resp.previousMessages };
        }
      });
      return obs;
    } else {
      // creating without proper authentication
      return undefined;
    }
  }

  join(roomName:string): Observable<JoinResponse> | undefined {
    const obs = this.roomBackendService.join(roomName);
    if (obs) {
      obs.subscribe({
        next: (resp:JoinResponse) => {
          this.roomAuthenticationInfo = { roomName: resp.roomName, previousMessages: resp.previousMessages };
        }
      });
      return obs;
    } else {
      // creating without proper authentication
      return undefined;
    }
  }

  delete(roomId:string): Observable<Object> {
    return this.roomBackendService.delete(roomId);
  }

  private roomAuthenticationInfo: RoomData | undefined = undefined;

}
