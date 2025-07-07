import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormPrediccionComponent } from './form-prediccion.component';

describe('FormPrediccionComponent', () => {
  let component: FormPrediccionComponent;
  let fixture: ComponentFixture<FormPrediccionComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [FormPrediccionComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(FormPrediccionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
