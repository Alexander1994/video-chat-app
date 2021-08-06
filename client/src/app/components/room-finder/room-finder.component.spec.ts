import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RoomFinderComponent } from './room-finder.component';

describe('RoomFinderComponent', () => {
  let component: RoomFinderComponent;
  let fixture: ComponentFixture<RoomFinderComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ RoomFinderComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(RoomFinderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
