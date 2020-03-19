import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ServerDetailComponent } from './server-detail.component';

describe('ServerDetailComponent', () => {
  let component: ServerDetailComponent;
  let fixture: ComponentFixture<ServerDetailComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ServerDetailComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ServerDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
