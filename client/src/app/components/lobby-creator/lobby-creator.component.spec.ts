import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LobbyCreatorComponent } from './lobby-creator.component';

describe('LobbyCreatorComponent', () => {
  let component: LobbyCreatorComponent;
  let fixture: ComponentFixture<LobbyCreatorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ LobbyCreatorComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(LobbyCreatorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
