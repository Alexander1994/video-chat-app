import { Component, OnInit } from '@angular/core';
import { FormBuilder, Validators } from '@angular/forms';
import { RoomService } from 'src/app/services/room-service/room.service'
import { Router } from '@angular/router';
import { JoinResponse } from 'src/app/services/roombackend-service/roombackend.service';

@Component({
  selector: 'app-room-finder',
  templateUrl: './room-finder.component.html',
  styleUrls: ['./room-finder.component.less']
})
export class RoomFinderComponent implements OnInit {

  roomFinderForm = this.formBuilder.group({
    roomName: ['', Validators.required],
  });

  constructor(private formBuilder: FormBuilder, private roomService:RoomService, private router:Router) {

  }

  onSubmit() {
    const roomName = this.roomFinderForm.value.roomName;
    const joinObs = this.roomService.join(roomName);
    if (joinObs) {
      joinObs.subscribe({
        next: (resp:JoinResponse) => {
          if (resp.success) {
            this.roomFinderForm.reset();
            console.log("RoomFinderComponent navigate /room");
  
            this.router.navigate(['/room']);
          } else {
            // display id failed
          }
        }
      });
    } else {
      // auth guard not working?
    }

  }

  ngOnInit(): void {

  }

}
